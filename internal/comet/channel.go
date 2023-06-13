package comet

import (
	"errors"
	"github.com/xyhubl/yim/pkg/bufio"
	"sync"

	"github.com/xyhubl/yim/api/protocol"
)

var ErrSignalFullMsgDropped = errors.New("channel: signal channel full, msg dropped")

type Channel struct {
	// zh: 连接IP
	IP string
	// zh: 连接唯一标识
	Key string
	// zh: 连接用户
	Mid int64
	// zh: 该连接所属的房间
	Room *Room

	Writer bufio.Writer
	Reader bufio.Reader

	Pre  *Channel
	Next *Channel

	CliProto Ring
	signal   chan *protocol.Proto

	watchOps map[int32]struct{}
	mutex    sync.RWMutex
}

func NewChannel(cli, svr int) *Channel {
	c := new(Channel)
	c.CliProto.Init(cli)

	c.signal = make(chan *protocol.Proto, svr)
	c.watchOps = make(map[int32]struct{})
	return c
}

// zh: 监听要发送的
func (c *Channel) Watch(accepts ...int32) {
	c.mutex.Lock()
	for _, op := range accepts {
		c.watchOps[op] = struct{}{}
	}
	c.mutex.Unlock()
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

func (c *Channel) Signal() {
	c.signal <- protocol.ProtoReady
}

func (c *Channel) Ready() *protocol.Proto {
	return <-c.signal
}

func (c *Channel) Close() {
	c.signal <- protocol.ProtoFinish
}
