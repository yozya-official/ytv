package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一返回结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Extra   interface{} `json:"extra,omitempty"`
}

// 成功返回
func Success(c *gin.Context, data interface{}, extra interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		Extra:   extra,
	})
}

// 错误返回
func Error(c *gin.Context, code int, message string, extra interface{}) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    code,
		Message: message,
		Extra:   extra,
	})
}
