package util

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type BaseResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func ReturnJson(code int, data interface{}, msg string, c *gin.Context) {
	result := &BaseResult{
		Code: code,
		Data: data,
		Msg:  msg,
	}
	jsonStr, _ := json.Marshal(result)
	c.Writer.WriteString(string(jsonStr))
}
