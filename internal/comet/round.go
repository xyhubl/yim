package comet

import (
	"github.com/xyhubl/yim/internal/comet/conf"
	"github.com/xyhubl/yim/pkg/bytes"
	"github.com/xyhubl/yim/pkg/time"
)

type RoundOptions struct {
	Timer        int
	TimerSize    int
	Reader       int
	ReadBuf      int
	ReadBufSize  int
	Writer       int
	WriteBuf     int
	WriteBufSize int
}

type Round struct {
	readers []bytes.Pool
	writers []bytes.Pool

	timers  []time.Timer
	options RoundOptions
}

func NewRound(c *conf.Config) (r *Round) {
	r = &Round{
		options: RoundOptions{
			Reader:       c.TCP.Reader,
			ReadBuf:      c.TCP.ReadBuf,
			ReadBufSize:  c.TCP.ReadBufSize,
			Writer:       c.TCP.Writer,
			WriteBuf:     c.TCP.WriteBuf,
			WriteBufSize: c.TCP.WriteBufSize,
			Timer:        c.Protocol.Timer,
			TimerSize:    c.Protocol.TimerSize,
		}}
	// zh: 初始化读取缓冲区
	// en: initialize read buffer
	r.readers = make([]bytes.Pool, r.options.Reader)
	for i := 0; i < r.options.Reader; i++ {
		r.readers[i].Init(r.options.ReadBuf, r.options.ReadBufSize)
	}
	// zh: 初始化写入缓冲区
	// en: initialize write buffer
	r.writers = make([]bytes.Pool, r.options.Writer)
	for i := 0; i < r.options.Writer; i++ {
		r.writers[i].Init(r.options.WriteBuf, r.options.WriteBufSize)
	}
	// zh: 初始化定时器数组
	// en: initialize timer array
	r.timers = make([]time.Timer, r.options.Timer)
	for i := 0; i < r.options.Timer; i++ {
		r.timers[i].Init(r.options.TimerSize)
	}
	return
}

func (r *Round) Timer(rn int) *time.Timer {
	return &(r.timers[rn%r.options.Timer])
}

func (r *Round) Reader(rn int) *bytes.Pool {
	return &(r.readers[rn%r.options.Reader])
}

func (r *Round) Writer(rn int) *bytes.Pool {
	return &(r.writers[rn%r.options.Writer])
}
