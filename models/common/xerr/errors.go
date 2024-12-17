package xerr

import (
	"fmt"
	"tz-gin/models/common/response"
)

type CustomError struct {
	Code    uint32 `json:"code"`
	Message string `json:"msg"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("errCode: %d,errMsg: %s", e.Code, e.Message)
}

func NewErrCodeMsg(errCode uint32, errMsg string) *CustomError {
	return &CustomError{Code: errCode, Message: errMsg}
}

// 根据错误码获取
func NewErrCode(errCode uint32) *CustomError {
	return &CustomError{
		Code:    errCode,
		Message: response.ErrorMessages[errCode],
	}
}

// 自定义错误
func NewErrMsg(errMsg string) *CustomError {
	return &CustomError{Code: response.SERVER_ERROR, Message: errMsg}
}

func UnknownError() *CustomError {
	return NewErrCode(response.SERVER_ERROR)
}

func ServerError() *CustomError {
	return NewErrCode(response.SERVER_ERROR)
}
