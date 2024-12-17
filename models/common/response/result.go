package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ERROR   = -1
	SUCCESS = 0
)

const (
	SERVER_ERROR    = 3 // 系统错误
	USER_NOT_PERMIT = 4 // 用户无权限
	PARAMETER_ERROR = 5 // 参数错误
	AUTH_ERROR      = 6 // 认证错误
	RESOURCE_EXIST  = 7 // 资源已存在
)

var ErrorMessages = map[uint32]string{
	SERVER_ERROR:    "系统错误",
	PARAMETER_ERROR: "参数错误",
	AUTH_ERROR:      "认证错误",
	RESOURCE_EXIST:  "资源已存在",
	USER_NOT_PERMIT: "暂无权限",
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

func Result(code int, data interface{}, success bool, message string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		message,
		success,
	})
}

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, true, "成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, true, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, true, "成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, true, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, false, "失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, false, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, false, message, c)
}

func FailWithErrorCode(code int, message string, c *gin.Context) {
	Result(code, map[string]interface{}{}, false, message, c)
}
