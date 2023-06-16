package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/xyhubl/yim/internal/logic/http/util"
)

func Ping(c *gin.Context) {
	util.ReturnJson(200, nil, "pong", c)
	return
}
