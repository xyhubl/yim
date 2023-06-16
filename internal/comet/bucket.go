package comet

import (
	"sync"

	"github.com/xyhubl/yim/internal/comet/conf"
)

type Bucket struct {
	c     *conf.Bucket
	cLock sync.RWMutex

	chs   map[string]*Channel
	rooms map[string]*Room

	ipCnts map[string]int32
}

func NewBucket(c *conf.Bucket) (b *Bucket) {
	b = new(Bucket)
	b.chs = make(map[string]*Channel, c.Channel)
	b.rooms = make(map[string]*Room, c.Room)

	b.ipCnts = make(map[string]int32)

	b.c = c

	return
}

func (b *Bucket) ChannelCount() int {
	return len(b.chs)
}

func (b *Bucket) RoomCount() int {
	return len(b.rooms)
}

func (b *Bucket) IPCount() (res map[string]struct{}) {
	b.cLock.RLock()
	res = make(map[string]struct{}, len(b.ipCnts))
	for ip := range b.ipCnts {
		res[ip] = struct{}{}
	}
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Put(rid string, ch *Channel) (err error) {
	b.cLock.Lock()
	// zh: 删除旧连接
	if oldCh := b.chs[ch.Key]; oldCh != nil {
		oldCh.Close()
	}
	b.chs[ch.Key] = ch
	var (
		ok   bool
		room *Room
	)
	if rid != "" {
		if room, ok = b.rooms[rid]; !ok {
			room = NewRoom(rid)
			b.rooms[rid] = room
		}
		ch.Room = room
	}
	// zh: 记录ip
	b.ipCnts[ch.IP]++
	b.cLock.Unlock()
	if room != nil {
		err = room.Put(ch)
	}
	return
}

func (b *Bucket) Del(ch *Channel) {
	room := ch.Room
	b.cLock.Lock()
	if oldCh, ok := b.chs[ch.Key]; ok {
		if oldCh == ch {
			delete(b.chs, oldCh.Key)
		}
		if b.ipCnts[oldCh.IP] > 1 {
			b.ipCnts[oldCh.IP]--
		} else {
			delete(b.ipCnts, oldCh.IP)
		}
	}
	b.cLock.Unlock()
	if room != nil && room.Del(ch) {
		// if empty room, must delete from bucket
		b.DelRoom(room)
	}
}

func (b *Bucket) Channel(key string) (ch *Channel) {
	b.cLock.RLock()
	ch = b.chs[key]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) Room(rid string) (room *Room) {
	b.cLock.RLock()
	room = b.rooms[rid]
	b.cLock.RUnlock()
	return
}

func (b *Bucket) DelRoom(room *Room) {
	b.cLock.Lock()
	delete(b.rooms, room.ID)
	b.cLock.Unlock()
	room.Close()
}

func (b *Bucket) Rooms() (res map[string]struct{}) {
	var (
		roomID string
		room   *Room
	)
	res = make(map[string]struct{})
	b.cLock.RLock()
	for roomID, room = range b.rooms {
		if room.Online > 0 {
			res[roomID] = struct{}{}
		}
	}
	b.cLock.RUnlock()
	return
}
