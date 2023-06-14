package comet

import (
	"github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/api/protocol"
	"github.com/xyhubl/yim/pkg/strings"
	"golang.org/x/net/context"
	"time"
)

// zh: 建立连接
func (s *Server) Connect(c context.Context, p *protocol.Proto, cookie string) (mid int64, key, rid string, accepts []int32, heartbeat time.Duration, err error) {
	reply, err := s.rpcClient.Connect(c, &logic.ConnectReq{
		Server: s.serverID,
		Cookie: cookie,
		Token:  p.Body,
	})
	if err != nil {
		return
	}
	return reply.Mid, reply.Key, reply.RoomID, reply.Accepts, time.Duration(reply.Heartbeat), nil
}

// zh: 关闭连接
func (s *Server) Disconnect(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Disconnect(c, &logic.DisconnectReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

// zh: 维持心跳
func (s *Server) Heartbeat(c context.Context, mid int64, key string) (err error) {
	_, err = s.rpcClient.Heartbeat(c, &logic.HeartbeatReq{
		Server: s.serverID,
		Mid:    mid,
		Key:    key,
	})
	return
}

func (s *Server) Operate(c context.Context, p *protocol.Proto, ch *Channel, b *Bucket) error {
	switch p.Op {
	// zh: 如果是改变房间
	case protocol.OpChangeRoom:
	// zh: 如果是订阅
	case protocol.OpSub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.Watch(ops...)
		}
		p.Op = protocol.OpSubReply
	// zh: 如果是取消订阅
	case protocol.OpUnsub:
		if ops, err := strings.SplitInt32s(string(p.Body), ","); err == nil {
			ch.UnWatch(ops...)
		}
		p.Op = protocol.OpUnsubReply
	}
	return nil
}
