package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"tz-gin/api"
	_ "tz-gin/docs"
	"tz-gin/global"
	"tz-gin/middleswares"
)

func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	global.DBClient = global.GormMysql()
	r.Use(middleswares.CorsMiddlewares())
	adminApi := api.AppApi.AdminApi
	adminRouter := r.Group("api/admin").Use(middleswares.AuthAdminCheck())
	{
		adminRouter.POST("/courses", middleswares.ErrWrapper(adminApi.AddCourse))
		adminRouter.DELETE("/courses/:courseId", middleswares.ErrWrapper(adminApi.DeleteCourse))
		adminRouter.PUT("/courses", middleswares.ErrWrapper(adminApi.UpdateCourse))
		adminRouter.GET("/courses", middleswares.ErrWrapper(adminApi.GetCourses))
		adminRouter.GET("/courses/:courseId", middleswares.ErrWrapper(adminApi.GetCoursesById))
		adminRouter.GET("/students", middleswares.ErrWrapper(adminApi.GetStudents))
		adminRouter.GET("/students/:studentId", middleswares.ErrWrapper(adminApi.GetStudent))
	}

	userApi := api.AppApi.UserApi
	userRouter := r.Group("api/user")
	userAuthRouter := r.Group("api/user").Use(middleswares.AuthUserCheck())
	{
		userRouter.POST("/register", middleswares.ErrWrapper(userApi.Register))
		userRouter.POST("", middleswares.ErrWrapper(userApi.Login))
		userAuthRouter.DELETE("", middleswares.ErrWrapper(userApi.Logout))
		userAuthRouter.GET("", middleswares.ErrWrapper(userApi.GetUser))
		userAuthRouter.GET("/courses", middleswares.ErrWrapper(userApi.GetCourses))
		userAuthRouter.GET("/courses/:courseId", middleswares.ErrWrapper(userApi.GetCourseById))
		userAuthRouter.POST("/courses", middleswares.ErrWrapper(userApi.EnrollCourse))
		userAuthRouter.DELETE("/courses/:courseId", middleswares.ErrWrapper(userApi.DropCourse))
		userAuthRouter.GET("/courses-selected", middleswares.ErrWrapper(userApi.GetEnrolledCourses))
		//userAuthRouter.GET("/schedule", middleswares.ErrWrapper(userApi.GetSchedule))
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	fmt.Println("swagger running on http://localhost:8282/swagger/index.html")
	r.Run(":8282")
	// Register
	// @Summary 学生注册接口
	// @Tags 学生部分
	// @Accept json
	// @Param body UserRegisterRequest true "学生注册请求参数"
	// @Success 200 {object} response.Response{} "success"
	// @Router /api/user/register [post]
}
