package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义了标准的API响应格式
type Response struct {
	ErrCode int         `json:"err_code"`
	Data    interface{} `json:"data,omitempty"`
	ErrMsg  string      `json:"err_msg,omitempty"`
}

// Success 响应成功的请求
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		ErrCode: 0,
		Data:    data,
	})
}

// Error 响应失败的请求
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, Response{
		ErrCode: code,
		ErrMsg:  msg,
	})
}
