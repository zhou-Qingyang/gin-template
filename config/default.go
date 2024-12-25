package config

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
)

func SetCORS(r *gin.Engine) {
	setConfig := cors.DefaultConfig()
	setConfig.AllowOrigins = split(Config.AllowOrigins)
	setConfig.AllowHeaders = split(Config.AllowHeaders)
	r.Use(cors.New(setConfig))
}

func split(s string) []string {
	return strings.Split(s, "|")
}
