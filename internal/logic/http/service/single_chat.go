package service

import (
	pb "github.com/xyhubl/yim/api/logic"
	"github.com/xyhubl/yim/internal/logic/dao"
	"github.com/xyhubl/yim/internal/logic/http/view_model"
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
		msg := pb.PushMsg{
			Type:      pb.Type_PUSH,
			Operation: req.Op,
			Server:    server,
			Keys:      v,
			Msg:       []byte(req.Msg),
		}
		if err = dao.PushMsg(ctx, req.Op, server, v, msg.Msg); err != nil {
			return err
		}
	}
	return nil
}
