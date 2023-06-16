package dao

import (
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xyhubl/yim/internal/logic/conf"
	kafka "gopkg.in/Shopify/sarama.v1"
)

var BaseDao *Base

type Base struct {
	c           *conf.Config
	Redis       *redis.Client
	RedisExpire time.Duration

	KafkaPub kafka.SyncProducer
}

func New(c *conf.Config) *Base {
	BaseDao = &Base{
		c:           c,
		KafkaPub:    newKafka(c.Kafka),
		Redis:       newRedis(c.Redis),
		RedisExpire: time.Duration(c.Redis.Expire) * time.Second,
	}
	return BaseDao
}

func newRedis(c *conf.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Network:      c.Network,
		Addr:         c.Addr,
		Username:     c.UserName,
		Password:     c.Password,
		DialTimeout:  time.Duration(c.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
	})
}

func newKafka(c *conf.Kafka) kafka.SyncProducer {
	kc := kafka.NewConfig()
	kc.Producer.RequiredAcks = kafka.WaitForAll // Wait for all in-sync replicas to ack the message
	kc.Producer.Retry.Max = 10                  // Retry up to 10 times to produce the message
	kc.Producer.Return.Successes = true
	pub, err := kafka.NewSyncProducer(c.Brokers, kc)
	if err != nil {
		panic(err)
	}
	return pub
}

func (b *Base) Close() error {
	return b.Redis.Close()
}
