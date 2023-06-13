package grpc

import (
	"net"
	"time"

	pb "github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/internal/logic"
	"github.com/xyhubl/yim/internal/logic/conf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func New(c *conf.RPCServer, l *logic.Logic) *grpc.Server {
	opts := grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle:     time.Second * time.Duration(c.IdleTimeout),       // 指定连接的最大空闲时间。如果连接在指定的空闲时间内没有出现任何活动（即没有进行 RPC 请求或连接建立），服务器将发送 GoAway 帧关闭连接。默认值是无限大，表示没有限制
		MaxConnectionAge:      time.Second * time.Duration(c.MaxLifeTime),       // 指定连接的最大存在时间。如果连接存在的时间超过指定的时间，服务器将发送 GoAway 帧关闭连接。默认值是无限大，表示没有限制。为了避免连接风暴，最大存在时间会加上一个随机抖动值
		MaxConnectionAgeGrace: time.Second * time.Duration(c.ForceCloseWait),    // 在 MaxConnectionAge 时间之后的额外宽限期。超过 MaxConnectionAge 时间后，服务器会强制关闭连接。默认值是无限大，表示没有额外宽限期
		Time:                  time.Second * time.Duration(c.KeepAliveInterval), // 在服务器没有收到任何活动后，通过发送 Ping 检查客户端是否仍然存活的时间间隔。如果设置的时间间隔小于 1 秒，将会使用最小值 1 秒。默认值是 2 小时。
		Timeout:               time.Second * time.Duration(c.KeepAliveTimeout),  // 在发送了 Ping 后，等待客户端活动的超时时间。如果在超时时间内没有收到任何活动，服务器将关闭连接。默认值是 20 秒。
	})
	srv := grpc.NewServer(opts)
	pb.RegisterLogicServer(srv, &server{
		srv: l,
	})
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
	pb.UnimplementedLogicServer
	srv *logic.Logic
}

var _ pb.LogicServer = &server{}

func (s *server) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectReply, error) {
	return nil, nil
}

func (s *server) Disconnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
	return nil, nil
}

func (s *server) Heartbeat(ctx context.Context, req *pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
	return nil, nil
}
