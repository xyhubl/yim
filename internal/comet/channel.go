package comet

import (
	"errors"
	"github.com/xyhubl/yim/api/protocol"
)

var (
	ErrSignalFullMsgDropped = errors.New("channel: signal channel full, msg dropped")
)

type Channel struct {
	// zh: 连接IP
	IP string
	// zh: 连接唯一标识
	Key string
	// zh: 连接用户
	Mid int64
	// zh: 该连接所属的房间
	Room *Room

	Pre  *Channel
	Next *Channel

	CliProto Ring
	signal   chan *protocol.Proto

	watchOps map[int32]struct{}
}

func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)

	c.signal = make(chan *protocol.Proto, svr)
	c.watchOps = make(map[int32]struct{})
	return c
}

func (c *Channel) Push(p *protocol.Proto) (err error) {
	select {
	case c.signal <- p:
		// zh: channel满了会进default
		// en: if channel is full, it will go to default
	default:
		err = ErrSignalFullMsgDropped
	}
	return
}

func (c *Channel) Ready() *protocol.Proto {
	return <-c.signal
}

func (c *Channel) Close() {
	c.signal <- protocol.ProtoFinish
}
