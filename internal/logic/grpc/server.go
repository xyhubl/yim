package grpc

import (
	pb "github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/internal/logic"
	"github.com/xyhubl/yim/internal/logic/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

func NewRpcSrv(c *conf.RPCServer, l *logic.Logic) *grpc.Server {
	srv := grpc.NewServer()
	pb.RegisterLogicServer(srv, &server{
		srv: l,
	})
	lis, err := net.Listen(c.Network, c.Addr)
	if err != nil {
		panic(err)
	}
	go func() {
		log.Println("[INFO] GRPC server start.", c.Addr)
		if err = srv.Serve(lis); err != nil {
			panic(err)
		}
	}()
	return srv
}

type server struct {
	pb.UnimplementedLogicServer
	srv *logic.Logic
}

var _ pb.LogicServer = &server{}

func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	mid, key, room, accepts, hb, err := s.srv.Connect(ctx, req.Server, req.Cookie, req.Token)
	if err != nil {
		return &pb.ConnectReply{}, err
	}
	return &pb.ConnectReply{Mid: mid, Key: key, RoomID: room, Accepts: accepts, Heartbeat: hb}, nil
}

func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	err := s.srv.DisConnect(ctx, req.Mid, req.Key, req.Server)
	return &pb.DisconnectReply{}, err
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	return nil, nil
}
