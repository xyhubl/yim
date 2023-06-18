package comet

import (
	"github.com/xyhubl/yim/api/logic"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/rand"
	"time"

	"github.com/xyhubl/yim/internal/comet/conf"
	"github.com/zhenjl/cityhash"
)

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
)

type Server struct {
	c         *conf.Config
	bucketIdx uint32
	serverID  string

	round     *Round
	buckets   []*Bucket
	rpcClient logic.LogicClient
}

func NewServer(c *conf.Config) *Server {
	s := &Server{
		c:         c,
		round:     NewRound(c),
		rpcClient: newLogicClient(c.RpcClient),
	}
	s.buckets = make([]*Bucket, c.Bucket.Size)
	s.bucketIdx = uint32(c.Bucket.Size)
	s.serverID = c.Base.Host
	for i := 0; i < c.Bucket.Size; i++ {
		s.buckets[i] = NewBucket(c.Bucket)
	}
	return s
}

// zh: 根据subKey将不同连接分配到不同桶
func (s *Server) Bucket(subKey string) *Bucket {
	idx := cityhash.CityHash32([]byte(subKey), uint32(len(subKey))) % s.bucketIdx
	return s.buckets[idx]
}

// zh: 随机心跳时间
func (s *Server) RandServerHeartbeat() time.Duration {
	return minServerHeartbeat + time.Duration(rand.Int63n(int64(maxServerHeartbeat-minServerHeartbeat)))
}

const (
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
)

// zh: logic rpc client
func newLogicClient(c *conf.RpcClient) logic.LogicClient {
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
	return logic.NewLogicClient(conn)
}
