package job

import (
	"errors"
	"github.com/google/martian/log"
	"github.com/xyhubl/yim/api/protocol"
	"github.com/xyhubl/yim/internal/job/conf"
	"github.com/xyhubl/yim/pkg/bytes"
	"time"
)

var (
	ErrRoomFull = errors.New("room: room proto chan full")
)

var (
	roomReadyProto = new(protocol.Proto)
)

type Room struct {
	c     *conf.Room
	job   *Job
	id    string
	proto chan *protocol.Proto
}

func NewRoom(job *Job, id string, c *conf.Room) (r *Room) {
	r = &Room{
		c:     c,
		id:    id,
		job:   job,
		proto: make(chan *protocol.Proto, c.Batch*2),
	}
	go r.pushProc(c.Batch, time.Millisecond*time.Duration(c.Signal))
	return
}

// zh: 获取房间
func (j *Job) getRoom(roomId string) *Room {
	j.roomsMtx.RLock()
	room, ok := j.rooms[roomId]
	j.roomsMtx.RUnlock()
	if !ok {
		j.roomsMtx.Lock()
		if room, ok = j.rooms[roomId]; !ok {
			room = NewRoom(j, roomId, j.c.Room)
			j.rooms[roomId] = room
		}
		j.roomsMtx.Unlock()
		log.Infof("[INFO] new room: %s, active room num: %d", roomId, len(j.rooms))
	}
	return room
}

// zh:删除房间
func (j *Job) delRoom(roomID string) {
	j.roomsMtx.Lock()
	delete(j.rooms, roomID)
	j.roomsMtx.Unlock()
}

// zh:room 发送消息逻辑
func (r *Room) Push(op int32, msg []byte) error {
	p := &protocol.Proto{
		Ver:  1,
		Op:   op,
		Body: msg,
	}
	var err error
	select {
	case r.proto <- p:
	default:
		err = ErrRoomFull
	}
	return err
}

// zh:room 发送消息等待逻辑
func (r *Room) pushProc(batch int, signalTime time.Duration) {
	log.Infof("[INFO] pushProc start roomId: %s", r.id)
	td := time.AfterFunc(signalTime, func() {
		select {
		case r.proto <- roomReadyProto:
		default:
		}
	})
	defer td.Stop()

	var (
		n    int
		last time.Time
		p    *protocol.Proto
		buf  = bytes.NewWriterSize(int(protocol.MaxBodySize))
	)
	for {
		// zh: 异常数据直接退出
		if p = <-r.proto; p == nil {
			break
		}
		// zh: 正常传输
		if p != roomReadyProto {
			p.WriteTo(buf)
			if n++; n == 1 {
				last = time.Now()
				td.Reset(signalTime)
				continue
				// zh: 如果还没达到数量要求 && 批量发送间隔时间还没有到,则进行下一次接收
			} else if n < batch {
				if signalTime > time.Since(last) {
					continue
				}
			}
		} else { // 长时间没有发送过消息，发送缓存的消息,销毁room 减少内存消耗
			if n == 0 {
				break
			}
		}
		// zh: 到这里说明消息满足条件 准备发送
		r.job.broadcastRoomRawBytes(r.id, buf.Buffer())
		// zh: 老的buf自动gc, 新开一个buffer
		buf = bytes.NewWriterSize(buf.Size())
		// zh: 重新计数
		n = 0
		// zh: 如果设置了空闲时间 则重新计时 否则默认1分钟
		if r.c.Idle != 0 {
			td.Reset(time.Duration(r.c.Idle) * time.Millisecond)
		} else {
			td.Reset(time.Minute)
		}
	}
	r.job.delRoom(r.id)
	log.Infof("[INFO] pushProc stop roomId: %s", r.id)
}
