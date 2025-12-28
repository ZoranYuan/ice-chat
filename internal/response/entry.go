package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
	{
		"code": 0,        // 业务状态码
		"msg":  "ok",     // 提示信息
		"data": {}        // 数据，可为 null
	}
*/

const (
	CodeOK           = 200
	CodeBadParam     = 400
	CodeUnauthorized = 401
	CodeServerErr    = 201
)

func write(ctx *gin.Context, httpCode int, bizCode int, msg string, data any) {
	ctx.JSON(httpCode, gin.H{
		"code": bizCode,
		"msg":  msg,
		"data": data,
	})
}

func OK(c *gin.Context) {
	write(c, http.StatusOK, CodeOK, "ok", nil)
}

func OKWithData(c *gin.Context, data any) {
	write(c, http.StatusOK, CodeOK, "ok", data)
}

func OKWithMsg(c *gin.Context, msg string) {
	write(c, http.StatusOK, CodeOK, msg, nil)
}

func BadRequestWithMessage(c *gin.Context, message string) {
	write(c, http.StatusBadRequest, CodeBadParam, message, nil)
}

func Unauthorized(c *gin.Context) {
	write(c, http.StatusUnauthorized, CodeUnauthorized, "unauthored", nil)
}

func InternalError(c *gin.Context, msg string) {
	write(c, http.StatusInternalServerError, CodeServerErr, msg, nil)
}

func Fail(c *gin.Context, bizCode int, msg string) {
	write(c, bizCode, bizCode, msg, nil)
}
