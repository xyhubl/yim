package logic

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xyhubl/yim/internal/logic/dao"
	"go.uber.org/zap"
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
		zap.L().Error("Connect Unmarshal err:" + err.Error())
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
	// 记录授权信息
	if mid > 0 {
		if err = dao.BaseDao.HSetExpire(c, dao.BaseDao.RedisExpire, dao.KeyMidServer(mid), key, server); err != nil {
			zap.L().Error("Connect HSetExpire err:" + err.Error())
			return
		}
		if err = dao.BaseDao.SetExpire(c, dao.KeyKeyServer(key), server, dao.BaseDao.RedisExpire); err != nil {
			zap.L().Error("Connect SetExpire err:" + err.Error())
			return
		}
	}
	return
}
