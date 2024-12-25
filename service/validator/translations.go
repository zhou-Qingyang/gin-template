package validator

import (
	ut "github.com/go-playground/universal-translator"
)

func timingTransZh(ut ut.Translator) error {
	return ut.Add("timing", "{0} 格式不正确，应该是 HH:mm 格式", false)
}
