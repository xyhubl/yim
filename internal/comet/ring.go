package comet

import (
	"errors"

	"github.com/xyhubl/yim/api/protocol"
)

var (
	ErrRingEmpty = errors.New("ring: buffer empty")
	ErrRingFull  = errors.New("ring: buffer full")
)

type Ring struct {
	num  uint64
	mask uint64

	rp   uint64
	wp   uint64
	data []protocol.Proto
}

func NewRing(num int) *Ring {
	r := new(Ring)
	r.init(uint64(num))
	return r
}

func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

func (r *Ring) init(num uint64) {
	// 必须是2^n
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= num - 1
		}
		num <<= 1
	}
	r.data = make([]protocol.Proto, num)
	r.num = num
	r.mask = r.num - 1
}

func (r *Ring) Get() (proto *protocol.Proto, err error) {
	if r.rp == r.wp {
		return nil, ErrRingEmpty
	}
	proto = &r.data[r.rp&r.mask]
	return
}

func (r *Ring) Set() (proto *protocol.Proto, err error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	proto = &r.data[r.wp&r.mask]
	return
}

func (r *Ring) GetAdv() {
	r.rp++
}

func (r *Ring) SetAdv() {
	r.wp++
}

func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
}
