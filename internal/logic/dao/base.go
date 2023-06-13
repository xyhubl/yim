package dao

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xyhubl/yim/internal/logic/conf"
)

type Base struct {
	c           *conf.Config
	redis       *redis.Client
	RedisExpire time.Duration
}

func New(c *conf.Config) *Base {
	return &Base{
		redis:       newRedis(c.Redis),
		RedisExpire: time.Duration(c.Redis.Expire) / time.Second,
	}
}

func newRedis(c *conf.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      c.Network,
		Addr:         c.Addr,
		DialTimeout:  time.Duration(c.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
	})
}

func (b *Base) Close() error {
	return b.redis.Close()
}
