package comet

import (
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

	round   *Round
	buckets []*Bucket
}

func NewServer(c *conf.Config) *Server {
	s := &Server{
		c:     c,
		round: NewRound(c),
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
