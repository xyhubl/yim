package grpc

import (
	"context"
	"errors"
	pb "github.com/xyhubl/yim/api/comet"
	"github.com/xyhubl/yim/internal/comet"
	"github.com/xyhubl/yim/internal/comet/conf"
	"google.golang.org/grpc"
	"net"
)

var (
	ErrPushMsgArg = errors.New("grpc: rpc push msg error")
)

func New(c *conf.RpcServer, s *comet.Server) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterCometServer(srv, &server{srv: s})
	lis, err := net.Listen(c.Network, c.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		if err = srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	pb.UnimplementedCometServer
	srv *comet.Server
}

func (s *server) PushMsg(ctx context.Context, req *pb.PushMsgReq) (reply *pb.PushMsgReply, err error) {
	if len(req.Keys) == 0 || req.Proto == nil {
		return nil, ErrPushMsgArg
	}
	for _, key := range req.Keys {
		bucket := s.srv.Bucket(key)
		if bucket == nil {
			continue
		}
		if channel := bucket.Channel(key); channel != nil {
			if !channel.NeedPush(req.ProtoOp) {
				continue
			}
			if err = channel.Push(req.Proto); err != nil {
				return
			}
		}
	}
	return &pb.PushMsgReply{}, nil
}
