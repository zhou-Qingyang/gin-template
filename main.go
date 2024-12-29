package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tz-gin/config"
	_ "tz-gin/docs"
	"tz-gin/router"
	"tz-gin/service/validator"
)

func main() {
	gin.SetMode(config.Config.AppMode)
	validator.InitValidator("zh")
	srv := router.NewServer()
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
		fmt.Printf("fail to init server: %s\n", err.Error())
	}
}
