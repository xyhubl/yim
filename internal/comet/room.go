package comet

import (
	"errors"
	"github.com/xyhubl/yim/api/protocol"
	"sync"
)

var (
	ErrRoomDroped = errors.New("room: droped")
)

type Room struct {
	// zh: 房间ID
	// en: room id
	ID    string
	rLock sync.RWMutex
	// zh: 连接链表
	// en: connection linked list
	next *Channel
	// zh: 是否丢弃
	// en: is drop
	drop bool
	// zh: 当前room的在线人数
	// en: current room online count
	Online int32
	// zh: 所有room的在线人数
	// en: all room online count
	AllOnline int32
}

func NewRoom(id string) (r *Room) {
	r = new(Room)
	r.ID = id
	r.drop = false
	r.next = nil
	r.Online = 0
	return
}

// zh: 将某个连接放入房间
// en: put channel into the room
func (r *Room) Put(ch *Channel) (err error) {
	r.rLock.Lock()
	if !r.drop {
		if r.next != nil {
			r.next.Pre = ch
		}
		ch.Next = r.next
		ch.Pre = nil
		r.next = ch
		r.Online++
	} else {
		err = ErrRoomDroped
	}
	r.rLock.Unlock()
	return
}

// zh: 删除房间中的某个连接
// en: delete channel from the room
func (r *Room) Del(ch *Channel) bool {
	r.rLock.Lock()
	if ch.Next != nil {
		ch.Next.Pre = ch.Pre
	}
	if ch.Pre != nil {
		ch.Pre.Next = ch.Next
	} else {
		r.next = ch.Next
	}
	ch.Next = nil
	ch.Pre = nil
	r.Online--
	r.drop = r.Online == 0
	r.rLock.Unlock()
	return r.drop
}

// zh: 对某个房间 发送消息
// en: push message to room
func (r *Room) Push(p *protocol.Proto) {
	r.rLock.Lock()
	for ch := r.next; ch != nil; ch = ch.Next {
		_ = ch.Push(p)
	}
	r.rLock.Unlock()
}

// zh: 关闭房间
// en: close room
func (r *Room) Close() {
	r.rLock.RLock()
	for ch := r.next; ch != nil; ch = ch.Next {
		ch.Close()
	}
	r.rLock.Unlock()
}

// zh: 所有房间在线数量
// en: all room online count
func (r *Room) OnlineNum() int32 {
	if r.AllOnline > 0 {
		return r.AllOnline
	}
	return r.Online
}
