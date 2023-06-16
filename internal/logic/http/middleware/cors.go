package middleware

import (
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
		c.Header("Access-Control-Allow-Methods", " GET, POST, PATCH, PUT, DELETE,OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		// 处理请求
		c.Next()
	}
}
