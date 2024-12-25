package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tz-gin/config"
	"tz-gin/controller"
)

func NewServer() *http.Server {
	r := gin.Default()
	config.SetCORS(r)
	//config.InitSession(r)
	InitRouter(r)
	s := &http.Server{
		Addr:    "0.0.0.0:8088",
		Handler: r,
	}
	return s
}

var ctr = new(controller.Controller)
