package core

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"example/internal/pkg/errno"
)

// Response 定义 API 响应结构
type Response struct {
	Code    int         `json:"code"`    // 错误码
	Message string      `json:"message"` // 错误消息
	Data    interface{} `json:"data"`    // 响应数据
}

// WriteResponse 封装了响应处理，统一返回 JSON 格式的数据
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		code, message := errno.Decode(err)
		c.JSON(http.StatusOK, Response{
			Code:    code,
			Message: message,
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    errno.OK.Code,
		Message: errno.OK.Message,
		Data:    data,
	})
}
