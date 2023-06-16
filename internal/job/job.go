package job

import (
	pb "github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/internal/job/conf"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"gopkg.in/Shopify/sarama.v1"
	"log"
)

type Job struct {
	c             *conf.Config
	cometServers  map[string]*Comet
	consumerGroup sarama.ConsumerGroup
}

func New(c *conf.Config) *Job {
	j := &Job{
		c: c,
	}
	// todo 由于还没有支持 服务发现注册 暂时写死
	j.cometServers = make(map[string]*Comet)
	j.cometServers["yim.comet"] = newComet(c)
	j.consumerGroup = newKafkaSub(c.Kafka)
	go func() {
		if err := j.consumerGroup.Consume(context.Background(), []string{c.Kafka.Topic}, j); err != nil {
			panic(err)
		}
	}()
	return j
}

func newKafkaSub(c *conf.Kafka) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_0_0_0 // 指定 Kafka 版本
	consumer, err := sarama.NewConsumerGroup(c.Brokers, c.Group, config)
	if err != nil {
		panic(err)
	}
	return consumer
}

func (j *Job) Close() error {
	if j.consumerGroup != nil {
		return j.consumerGroup.Close()
	}
	return nil
}

// 定义消费者组处理函数
func (j *Job) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (j *Job) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (j *Job) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		pushMsg := new(pb.PushMsg)
		if err := proto.Unmarshal(msg.Value, pushMsg); err != nil {
			log.Println("[ERROR] consumer Unmarshal err: ", err, "msg: ", msg)
			continue
		}
		if err := j.push(context.Background(), pushMsg); err != nil {
			log.Println("[ERROR] consumer push err: ", err)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}
