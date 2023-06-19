package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/xyhubl/yim/internal/logic/http/controller"
)

func InitRouter(middlewares ...gin.HandlerFunc) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares...)
	router.GET("/ping", controller.Ping)

	v1 := router.Group("v1")
	{
		push := v1.Group("/push")
		push.POST("/keys", controller.PushKeys)
		push.POST("/mids", controller.PushMids)
		push.POST("/room", controller.PushRoom)
	}
	return router
}
