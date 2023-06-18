package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xyhubl/yim/internal/logic/http/service"
	"github.com/xyhubl/yim/internal/logic/http/util"
	"github.com/xyhubl/yim/internal/logic/http/view_model"
	"golang.org/x/net/context"
)

func PushKeys(c *gin.Context) {
	req := new(view_model.PushKeysReq)
	if err := util.ValidParams(c, req); err != nil {
		util.ReturnJson(400, nil, "参数错误", c)
		return
	}
	ctx := context.WithValue(context.Background(), "trace_id", c.Value("trace_id"))
	if err := service.PushKeys(ctx, req); err != nil {
		util.ReturnJson(400, nil, fmt.Sprintf(" pushKeys err: %v", err), c)
	}
	util.ReturnJson(200, nil, "成功", c)
	return
}

func PushMids(c *gin.Context) {
	req := new(view_model.PushMidsReq)
	if err := util.ValidParams(c, req); err != nil {
		util.ReturnJson(400, nil, "参数错误", c)
		return
	}
	ctx := context.WithValue(context.Background(), "trace_id", c.Value("trace_id"))
	if err := service.PushMids(ctx, req); err != nil {
		util.ReturnJson(400, nil, fmt.Sprintf(" pushKeys err: %v", err), c)
	}
	util.ReturnJson(200, nil, "成功", c)
	return
}
