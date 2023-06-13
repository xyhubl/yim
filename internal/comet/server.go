package comet

import "github.com/xyhubl/yim/internal/comet/conf"

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
