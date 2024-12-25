package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"tz-gin/middleswares"
)

func InitRouter(r *gin.Engine) {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Use(middleswares.GinLogger(), middleswares.GinRecovery(true))

	adminApi := ctr.AdminApi
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

	userApi := ctr.UserApi
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
	r.POST("/test", middleswares.ErrWrapper(userApi.TestValidate))

	fmt.Println("swagger running on http://localhost:8282/swagger/index.html")
}
