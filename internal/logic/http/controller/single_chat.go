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
	return
}
