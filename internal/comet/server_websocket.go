package comet

import (
	"io"
	"log"
	"math"
	"net"
	"strings"
	"time"

	"github.com/xyhubl/yim/api/protocol"
	"github.com/xyhubl/yim/pkg/bytes"
	xtime "github.com/xyhubl/yim/pkg/time"
	"github.com/xyhubl/yim/pkg/websocket"
	"golang.org/x/net/context"
)

func InitWebsocket(server *Server, addrs []string, accept int) (err error) {
	var (
		bind     string
		addr     *net.TCPAddr
		listener *net.TCPListener
	)
	for _, bind = range addrs {
		if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
			return
		}
		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			return
		}
		log.Println("[INFO] websocket listen ", bind)
		for i := 0; i < accept; i++ {
			go acceptWebsocket(server, listener)
		}
	}
	return
}

func acceptWebsocket(server *Server, lis *net.TCPListener) {
	var (
		r    int
		err  error
		conn *net.TCPConn
	)
	for {
		if conn, err = lis.AcceptTCP(); err != nil {
			log.Printf("[ERROR]: AcceptTCP err: %v", err)
			return
		}
		if err = conn.SetKeepAlive(server.c.TCP.KeepAlive); err != nil {
			log.Printf("[ERROR]: SetKeepAlive err: %v", err)
			return
		}
		if err = conn.SetReadBuffer(server.c.TCP.RcvBuf); err != nil {
			log.Printf("[ERROR]: SetReadBuffer err: %v", err)
			return
		}
		if err = conn.SetWriteBuffer(server.c.TCP.SndBuf); err != nil {
			log.Printf("[ERROR]: SetWriteBuffer err: %v", err)
			return
		}
		go server.ServeWebsocket(conn, server.round.Reader(r), server.round.Writer(r), server.round.Timer(r))
		r++
		if r == math.MaxInt {
			r = 0
		}
	}
}

func (s *Server) ServeWebsocket(conn net.Conn, rp, wp *bytes.Pool, tr *xtime.Timer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		err error
		// zh: 最后一次心跳
		lastHb = time.Now()
		ws     *websocket.Conn
		req    *websocket.Request
	)
	step := 0
	// zh: 添加握手的超时机制
	trd := tr.Add(time.Duration(s.c.Protocol.HandshakeTimeout)*time.Second, func() {
		_ = conn.SetDeadline(time.Now().Add(time.Millisecond * 100))
		_ = conn.Close()
		log.Printf("[ERROR]: handshake timeout remoteIP: %s, step: %d", conn.RemoteAddr(), step)
	})

	// zh: 获取IP
	ch := NewChannel(s.c.Protocol.CliProto, s.c.Protocol.SvrProto)
	ch.IP, _, _ = net.SplitHostPort(conn.RemoteAddr().String())

	step = 1
	// zh: 读取http头部信息
	rb := rp.Get()
	ch.Reader.ResetBuffer(conn, rb.Bytes())
	rr := &ch.Reader
	if req, err = websocket.ReadRequest(rr); err != nil || req.RequestURI != "/sub" {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		if err != io.EOF {
			log.Printf("[ERROR]: websocket ReadRequest err: %v, step: %d", err, step)
		}
		return
	}

	step = 2
	// zh: 升级成websocket协议
	wb := wp.Get()
	ch.Writer.ResetBuffer(conn, wb.Bytes())
	rw := &ch.Writer
	if ws, err = websocket.Upgrade(conn, rr, rw, req); err != nil {
		conn.Close()
		tr.Del(trd)
		rp.Put(rb)
		wp.Put(wb)
		if err != io.EOF {
			log.Printf("[ERROR]: websocket Upgrade err: %v, step: %d", err, step)
		}
		return
	}

	step = 3
	// zh: 授权认证
	var (
		rid     string
		accepts []int32
		p       *protocol.Proto
		hb      time.Duration
		b       *Bucket
	)
	if p, err = ch.CliProto.Set(); err == nil {
		if ch.Mid, ch.Key, rid, accepts, hb, err = s.authWebsocket(ctx, ws, p, req.Header.Get("Cookie")); err == nil {
			// zh: 监听需要的连接
			ch.Watch(accepts...)
			// zh: 根据hash分配到不同的桶中
			b = s.Bucket(ch.Key)
			// zh: 将连接放入到room
			err = b.Put(rid, ch)
		}
	}
	if err != nil {
		ws.Close()
		rp.Put(rb)
		wp.Put(wb)
		tr.Del(trd)
		if err != io.EOF && err != websocket.ErrMessageClose {
			log.Printf("[ERROR]: websocket authWebsocket err: %v, key: %s, remoteIP: %s, step: %d", err, ch.Key, conn.RemoteAddr().String(), step)
		}
		return
	}
	step = 4
	// zh: 持续监听事件
	go s.dispatchWebsocket(ws, wp, wb, ch)
	serverHeartbeat := s.RandServerHeartbeat()
	for {
		if p, err = ch.CliProto.Set(); err != nil {
			break
		}
		if err = p.ReadWebsocket(ws); err != nil {
			break
		}
		if p.Op == protocol.OpHeartbeat {
			// zh: 如果是心跳,则重新计时
			tr.Set(trd, hb)
			p.Op = protocol.OpHeartbeatReply
			p.Body = nil
			if now := time.Now(); now.Sub(lastHb) > serverHeartbeat {
				if errHeartBeat := s.Heartbeat(ctx, ch.Mid, ch.Key); errHeartBeat == nil {
					lastHb = now
				} else {
					// zh: 错误的话,不用做处理,定时器会自动执行,释放内存
					log.Printf("[ERROR] ServeWebsocket Heartbeat err" + errHeartBeat.Error())
				}
			}
			step++
		} else {
			// zh: 如果不是心跳事件
			if err = s.Operate(ctx, p, ch, b); err != nil {
				break
			}
		}
		// zh: 写计数+1
		ch.CliProto.SetAdv()
		ch.Signal()
	}
	b.Del(ch)
	tr.Del(trd)
	ws.Close()
	ch.Close()
	rp.Put(rb)
	// zh: 意外的关闭
	if err != nil && err != io.EOF && err != websocket.ErrMessageClose && !strings.Contains(err.Error(), "closed") {
		log.Printf("[ERROR]: The closing of the accident key: %s, err: %v ", ch.Key, err)
	}
	// zh: 关闭连接
	if err = s.Disconnect(ctx, ch.Mid, ch.Key); err != nil {
		log.Printf("[ERROR]: Disconnect err mid: %d, key: %s, err: %v ", ch.Mid, ch.Key, err)
	}
}

func (s *Server) authWebsocket(ctx context.Context, ws *websocket.Conn, p *protocol.Proto, cookie string) (mid int64, key, rid string,
	accepts []int32, hb time.Duration, err error) {
	for {
		if err = p.ReadWebsocket(ws); err != nil {
			return
		}
		if p.Op == protocol.OpAuth {
			break
		} else {
			log.Printf("[ERROR] request not auth, op: %d", p.Op)
		}
	}
	if mid, key, rid, accepts, hb, err = s.Connect(ctx, p, cookie); err != nil {
		return
	}
	p.Op = protocol.OpAuthReply
	p.Body = nil
	if err = p.WriteWebsocket(ws); err != nil {
		return
	}
	err = ws.Flush()
	return
}

func (s *Server) dispatchWebsocket(ws *websocket.Conn, wp *bytes.Pool, wb *bytes.Buffer, ch *Channel) {
	var (
		err    error
		finish bool
		online int32
	)
	for {
		// zh: 阻塞直到有消息进来
		var p = ch.Ready()
		switch p {
		case protocol.ProtoFinish:
			finish = true
			goto failed
		// -------------------------------------------------------------------------------------------------------------
		case protocol.ProtoReady:
			for {
				// zh: 如果获取不到 说明写完了 或者 一般是队列满了拿不到,存在刷消息的嫌疑,直接丢弃
				if p, err = ch.CliProto.Get(); err != nil {
					break
				}
				// zh: 如果是心跳
				if p.Op == protocol.OpHeartbeatReply {
					if ch.Room != nil {
						online = ch.Room.OnlineNum()
					}
					// zh: 发送心跳回应
					if err = p.WriteWebsocketHeart(ws, online); err != nil {
						goto failed
					}
				} else {
					if err = p.WriteWebsocket(ws); err != nil {
						goto failed
					}
				}
				// zh: 避免内存泄漏
				p.Body = nil
				ch.CliProto.GetAdv()
			}
		// -------------------------------------------------------------------------------------------------------------
		default:
			if err = p.WriteWebsocket(ws); err != nil {
				goto failed
			}
		}
		// zh: hungry flush
		if err = ws.Flush(); err != nil {
			break
		}
	}
failed:
	ws.Close()
	wp.Put(wb)
	// zh: 丢弃还未读的消息, 不然reader将会堵塞
	for !finish {
		finish = (ch.Ready() == protocol.ProtoFinish)
	}
}
