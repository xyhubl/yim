package service

import (
	"fmt"
	"github.com/xyhubl/yim/internal/logic/dao"
	"github.com/xyhubl/yim/internal/logic/http/view_model"
	"golang.org/x/net/context"
	"net/url"
)

func EncodeRoomKey(typ string, room string) string {
	return fmt.Sprintf("%s://%s", typ, room)
}

func DecodeRoomKey(key string) (string, string, error) {
	u, err := url.Parse(key)
	if err != nil {
		return "", "", err
	}
	return u.Scheme, u.Host, nil
}

func PushRoom(ctx context.Context, req *view_model.PushRoomReq) error {
	return dao.BroadcastRoomMsg(ctx, req.Op, EncodeRoomKey(req.Typ, req.Room), []byte(req.Body))
}
