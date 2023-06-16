package job

import (
	"fmt"
	"github.com/xyhubl/yim/api/comet"
	"github.com/xyhubl/yim/internal/job/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync/atomic"
	"time"
)

var (
	// grpc options
	grpcMaxSendMsgSize = 1 << 24
	grpcMaxCallMsgSize = 1 << 24
)

const (
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
)

func newCometClient(c *conf.RpcClient) comet.CometClient {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Dial)*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, c.Addr, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithInitialWindowSize(grpcInitialWindowSize),
		grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize), grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	}...)
	if err != nil {
		panic(err)
	}
	return comet.NewCometClient(conn)
}

type Comet struct {
	serverID string
	client   comet.CometClient

	pushChanNum uint64                   // // 单聊消息数量
	pushChan    []chan *comet.PushMsgReq // 单聊消息

	routineSize uint64

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func newComet(c *conf.Config) *Comet {
	cmt := &Comet{
		serverID:    "yim.comet", // 暂定一个
		pushChan:    make([]chan *comet.PushMsgReq, c.Comet.RoutineSize),
		routineSize: uint64(c.Comet.RoutineSize),
	}
	cmt.client = newCometClient(c.RpcClient)
	cmt.ctx, cmt.ctxCancel = context.WithCancel(context.Background())

	for i := 0; i < c.Comet.RoutineSize; i++ {
		cmt.pushChan[i] = make(chan *comet.PushMsgReq, c.Comet.RoutineChan)
		go cmt.process(cmt.pushChan[i])
	}
	return cmt
}

func (c *Comet) Push(req *comet.PushMsgReq) {
	idx := atomic.AddUint64(&c.pushChanNum, 1) % c.routineSize
	c.pushChan[idx] <- req
	return
}

func (c *Comet) process(pushChan chan *comet.PushMsgReq) {
	for {
		select {
		case req := <-pushChan:
			if _, err := c.client.PushMsg(context.Background(), &comet.PushMsgReq{
				Keys:    req.Keys,
				Proto:   req.Proto,
				ProtoOp: req.ProtoOp,
			}); err != nil {
				log.Println("[ERROR] process PushMsg err: ", err)
			}
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *Comet) Close() error {
	var (
		err    error
		finish = make(chan bool)
	)
	go func() {
		for {
			n := 0
			for _, ch := range c.pushChan {
				n += len(ch)
			}
			if n == 0 {
				finish <- true
				return
			}
		}
	}()
	select {
	case <-finish:
		log.Println("[INFO] comet finish")
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("close comet(server:%s push: %d ) timeout", c.serverID, len(c.pushChan))
	}
	c.ctxCancel()
	return err
}
