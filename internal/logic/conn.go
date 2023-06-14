package logic

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xyhubl/yim/internal/logic/dao"
	"golang.org/x/net/context"
	"time"
)

type connectParam struct {
	Mid      int64   `json:"mid"`
	Key      string  `json:"key"`
	RoomID   string  `json:"room_id"`
	Platform string  `json:"platform"`
	Accepts  []int32 `json:"accepts"`
}

func (l *Logic) Connect(c context.Context, server, cookie string, token []byte) (mid int64, key, roomID string, accepts []int32, hb int64, err error) {
	param := new(connectParam)
	if err = json.Unmarshal(token, param); err != nil {
		return
	}
	mid = param.Mid
	roomID = param.RoomID
	accepts = param.Accepts
	hb = int64(8 * time.Minute)
	key = param.Key
	if key == "" {
		key = uuid.New().String()
	}
	if mid > 0 {
		if err = l.DaoBase.HSetExpire(c, l.DaoBase.RedisExpire, dao.KeyMidServer(mid), key, server); err != nil {
			return
		}
		if err = l.DaoBase.SetExpire(c, dao.KeyKeyServer(key), server, l.DaoBase.RedisExpire); err != nil {
			return
		}
	}
	return
}
