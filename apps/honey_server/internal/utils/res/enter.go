package res

import (
	"github.com/gin-gonic/gin"
	"honey_app/apps/honey_server/internal/utils/validate"
)

type Response struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func response(code int, data any, msg string, c *gin.Context) {
	c.JSON(200, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

func Ok(data any, msg string, c *gin.Context) {
	response(0, data, msg, c)
}

func OkWithData(data any, c *gin.Context) {
	Ok(data, "ok", c)
}

func OkWithMsg(msg string, c *gin.Context) {
	Ok(gin.H{}, msg, c)
}

func OkWithList(list any, count int64, c *gin.Context) {
	Ok(gin.H{"list": list, "count": count}, "成功", c)
}

func Fail(code int, msg string, c *gin.Context) {
	response(code, nil, msg, c)
}

func FailWithMsg(msg string, c *gin.Context) {
	response(1001, nil, msg, c)
}

func FailWithError(err error, c *gin.Context) {
	errMsg := validate.ValidateError(err)
	response(1001, nil, errMsg, c)
}
