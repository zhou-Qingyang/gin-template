package api

type Api struct {
	UserApi  UserApi
	AdminApi AdminApi
}

var AppApi = new(Api)
