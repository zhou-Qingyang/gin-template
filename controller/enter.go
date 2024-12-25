package controller

import "tz-gin/service"

type Controller struct {
	UserApi  UserApi
	AdminApi AdminApi
}

var services = new(service.Service)
