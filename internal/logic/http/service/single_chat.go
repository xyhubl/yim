package service

import (
	"github.com/xyhubl/yim/internal/logic/dao"
	"github.com/xyhubl/yim/internal/logic/http/view_model"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func PushKeys(ctx context.Context, req *view_model.PushKeysReq) error {
	keyList := make([]string, 0, len(req.Keys))
	for _, v := range req.Keys {
		keyList = append(keyList, dao.KeyKeyServer(v))
	}
	servers, err := dao.BaseDao.MGetString(ctx, keyList)
	if err != nil {
		return err
	}
	pushKeys := make(map[string][]string)
	for i, v := range req.Keys {
		server := servers[i]
		if server != "" && v != "" {
			pushKeys[server] = append(pushKeys[server], v)
		}
	}
	// 发送消息
	for server, v := range pushKeys {
		if err = dao.PushMsg(ctx, req.Op, server, v, []byte(req.Msg)); err != nil {
			return err
		}
	}
	return nil
}

func PushMids(ctx context.Context, req *view_model.PushMidsReq) error {
	keysAndServerMap := make(map[string][]string)
	for _, v := range req.Mids {
		resMap, err := dao.BaseDao.HGetAll(ctx, dao.KeyMidServer(v))
		if err != nil {
			zap.L().Error("PushMids HGetAll err: " + err.Error())
			return err
		}
		for key, server := range resMap {
			keysAndServerMap[server] = append(keysAndServerMap[server], key)
		}
	}
	for server, keys := range keysAndServerMap {
		if err := dao.PushMsg(ctx, req.Op, server, keys, []byte(req.Msg)); err != nil {
			zap.L().Error("PushMids PushMsg err: " + err.Error())
			return err
		}
	}
	return nil
}
