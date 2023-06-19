package controller

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xyhubl/yim/internal/logic/http/service"
	"github.com/xyhubl/yim/internal/logic/http/util"
	"github.com/xyhubl/yim/internal/logic/http/view_model"
)

func PushRoom(c *gin.Context) {
	req := new(view_model.PushRoomReq)
	if err := util.ValidParams(c, req); err != nil {
		util.ReturnJson(400, nil, "参数错误", c)
		return
	}
	ctx := context.WithValue(context.Background(), "trace_id", c.Value("trace_id"))
	if err := service.PushRoom(ctx, req); err != nil {
		util.ReturnJson(400, nil, fmt.Sprintf(" PushRoom err: %v", err), c)
		return
	}
	util.ReturnJson(200, nil, "成功", c)
	return
}
