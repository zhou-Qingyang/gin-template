package validator

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

// 校验函数
func timing(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	fmt.Println("validator1", fieldValue)
	if fieldValue == "" {
		// 处理空字符串的情况，可能是设为默认时间，或者报错等
		return false // 返回一个零值时间或者默认时间
	}
	// 尝试将字符串解析为时间
	parsedTime, err := time.Parse("2006-01-02 15:04:05", fieldValue)
	if err != nil {
		fmt.Println("validator2", err)
		return false // 如果解析失败，验证失败
	}
	fmt.Println("validator3", parsedTime)
	today := time.Now()
	if today.After(parsedTime) {
		return false // 如果时间早于今天，验证失败
	}
	return true // 验证通过
}
