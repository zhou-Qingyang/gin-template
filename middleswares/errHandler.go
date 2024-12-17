package middleswares

import (
	"tz-gin/models/common/response"
	"tz-gin/models/common/xerr"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error

func ErrWrapper(handler HandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		var (
			err error
		)
		err = handler(c)
		if err != nil {
			var apiException *xerr.CustomError
			if h, ok := err.(*xerr.CustomError); ok {
				apiException = h
			} else if _, ok := err.(error); ok {
				apiException = xerr.UnknownError()
			} else {
				apiException = xerr.ServerError()
			}
			response.FailWithErrorCode(int(apiException.Code), apiException.Message, c)
			return
		}
	}
}
