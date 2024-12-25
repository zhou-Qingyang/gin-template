package utils

import (
	"github.com/go-playground/validator/v10"
	validator2 "tz-gin/service/validator"
)

func ParseErr(err validator.ValidationErrors) string {
	ms := err.Translate(validator2.Trans)
	message := []string{}
	if len(ms) == 0 {
		return "Invalid input"
	}
	for _, v := range ms {
		message = append(message, v)
	}
	return message[0]
}
