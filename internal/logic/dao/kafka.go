package dao

import (
	pb "github.com/xyhubl/yim/api/logic"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
	"gopkg.in/Shopify/sarama.v1"
)

func PushMsg(ctx context.Context, op int32, server string, keys []string, msg []byte) error {
	pushMsg := &pb.PushMsg{
		Type:      pb.Type_PUSH,
		Operation: op,
		Server:    server,
		Keys:      keys,
		Msg:       msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return err
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(keys[0]),
		Topic: BaseDao.c.Kafka.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = BaseDao.KafkaPub.SendMessage(m); err != nil {
		return err
	}
	return nil
}

func BroadcastRoomMsg(c context.Context, op int32, room string, msg []byte) error {
	pushMsg := &pb.PushMsg{
		Type:      pb.Type_ROOM,
		Operation: op,
		Room:      room,
		Msg:       msg,
	}
	b, err := proto.Marshal(pushMsg)
	if err != nil {
		return err
	}
	m := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(room),
		Topic: BaseDao.c.Kafka.Topic,
		Value: sarama.ByteEncoder(b),
	}
	if _, _, err = BaseDao.KafkaPub.SendMessage(m); err != nil {
		return err
	}
	return nil
}
