package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
	"time"
	"tz-gin/global"
	"tz-gin/models"
	"tz-gin/models/common/response"
	"tz-gin/models/common/xerr"
	"tz-gin/utils"
)

type UserApi struct{}

type UserRegisterRequest struct {
	Name      string `json:"name"`
	Password  string `json:"password"`
	StudentId string `json:"studentId"`
}

// Register
// @Summary 学生注册接口
// @Description 注册接口
// @Tags 学生部分
// @Accept json
// @Param req body UserRegisterRequest true "学生注册请求参数"
// @Success 200 {object} response.Response{} "success"
// @Router /api/user/register [post]
func (u *UserApi) Register(c *gin.Context) error {
	var createUser UserRegisterRequest
	if err := c.ShouldBindJSON(&createUser); err != nil {
		return err
	}
	if createUser.Name == "" || createUser.Password == "" || createUser.StudentId == "" {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	userInfo, err := services.UserService.FindByStudentId(createUser.StudentId)
	if err != nil {
		return err
	}
	if userInfo.Id > 0 {
		return xerr.NewErrCode(response.RESOURCE_EXIST)
	}
	if err := services.UserService.CreateAccount(&models.UserModel{
		StudentName: createUser.Name,
		Password:    utils.Md5(createUser.Password),
		StudentId:   createUser.StudentId,
		Id:          userInfo.Id,
	}).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	response.Ok(c)
	return nil
}

type UserLoginRequest struct {
	StudentId string `json:"studentId"`
	Password  string `json:"password"`
}

// Login
// @Summary 学生登录接口
// @Tags 学生部分
// @Accept json
// @Param req body UserLoginRequest true "学生登录请求参数"
// @Success 200 {object} response.Response{} "success"
// @Router /api/user/ [post]
func (u *UserApi) Login(c *gin.Context) error {
	var loginUser UserLoginRequest
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		return err
	}
	if loginUser.StudentId == "" || loginUser.Password == "" {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	user, err := services.UserService.FindByStudentId(loginUser.StudentId)
	if err != nil {
		return err
	}

	if user.Password != utils.Md5(loginUser.Password) {
		return xerr.NewErrCodeMsg(3, "密码错误")
	}

	var student models.StudentModel
	if err := models.DB.Model(&models.StudentModel{}).
		Where("user_id = ?", user.Id).
		First(&student).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	// 返回token
	jwtUtils := utils.NewJWT()
	customClaims := jwtUtils.CreateClaims(utils.BaseClaims{
		StudentId:        user.StudentId,
		StudentName:      user.StudentName,
		UserId:           user.Id,
		StudentPrimaryId: student.Id,
		IsAdmin:          user.IsAdmin,
	})

	token, err := jwtUtils.CreateToken(customClaims)
	if err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	// 设置本地存储
	global.LocalCache.Set(fmt.Sprintf("token_%s", token), user.Id, 1*time.Hour)
	// 获取
	cachedCourses, found := global.LocalCache.Get(fmt.Sprintf("token_%s", token))
	if found {
		fmt.Println("用户id", cachedCourses)
	} else {
		response.Fail(c)
		return nil
	}
	response.OkWithData(map[string]interface{}{
		"token": token,
	}, c)
	return nil
}

type UserGetCourseRequest struct {
	StudentId   string `json:"studentId"`
	StudentName string `json:"studentName"`
}

// GetUser
// @Summary 查询学生信息
// @Tags 学生部分
// @Success 200 {object} response.Response{data=UserGetCourseRequest} "success"
// @Router /api/user [get]
func (u *UserApi) GetUser(c *gin.Context) error {
	baseClaims, exist := c.Get("claims")
	if !exist {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	claims := baseClaims.(*utils.CustomClaims)
	user, err := services.UserService.FindUserByUserId(claims.UserId)
	if err != nil {
		return err
	}
	response.OkWithData(UserGetCourseRequest{
		StudentId:   user.StudentId,
		StudentName: user.StudentName,
	}, c)
	return nil
}

type UserGetCoursesRequestTo struct {
	CourseName string   `form:"courseName"`
	Page       int      `form:"page"`
	Limit      int      `form:"limit"`
	Location   string   `form:"location"`
	Teachers   []string `form:"teachers"`
}

// GetCourses
// @Summary 查询课程列表
// @Description 查询课程列表
// @Accept json
// @Param req query UserGetCoursesRequestTo true "查询课程列表请求参数"
// @Tags 学生部分
// @Success 200 {object} response.Response{data=CoursesResponse} "success"
// @Router /api/user/courses [get]
func (u *UserApi) GetCourses(c *gin.Context) error {
	var req UserGetCoursesRequestTo
	if err := c.BindQuery(&req); err != nil {
		return err
	}
	count, tmp, err := services.UserService.ListCoursesBy(req.CourseName, req.Location, req.Page, req.Limit)
	if err != nil {
		return err
	}
	var courseTeacherMap = make(map[int64][]string)
	if len(tmp) > 0 {
		teacherDb := models.DB.Model(&tmp)

		tmpCourses := make([]models.CourseModel, 0)
		if err := teacherDb.
			Preload("Teachers").
			Find(&tmpCourses).Error; err != nil {
			return xerr.NewErrCode(response.SERVER_ERROR)
		}

		for _, course := range tmpCourses {
			for _, teacher := range course.Teachers {
				courseTeacherMap[course.Id] = append(courseTeacherMap[course.Id], teacher.TeacherName)
			}
		}
	}

	var rows []CourseItem
	for _, course := range tmp {
		var times []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}
		times = append(times, struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}{
			StartTime: course.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:   course.EndTime.Time.Format("2006-01-02 15:04:05"),
		})
		rows = append(rows, CourseItem{
			Id:         course.Id,
			CourseName: course.CourseName,
			Capacity:   course.Capacity,
			Location:   course.Location,
			Time:       times,
			Teachers:   courseTeacherMap[course.Id],
		})
	}

	response.OkWithData(CoursesResponse{
		Size: int(count),
		Rows: rows,
	}, c)
	return nil
}

// GetCourseById
// @Summary 查询课程信息
// @Tags 学生部分
// @Accept json
// @Param courseId path int true "课程id"
// @Success 200 {object} response.Response{data=GetCoursesByIdResponse} "success"
// @Router /api/user/courses/{courseId} [get]
func (u *UserApi) GetCourseById(c *gin.Context) error {
	courseIdParam := c.Param("courseId")
	courseId, err := strconv.Atoi(courseIdParam)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	course, err := services.UserService.FindCourseByCourseId(courseId)
	if err != nil {
		return err
	}

	var res GetCoursesByIdResponse
	res.Id = course.Id
	res.CourseName = course.CourseName
	res.Capacity = course.Capacity
	res.Location = course.Location
	res.Time = append(res.Time, struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	}{
		StartTime: course.StartTime.Time.Format("2006-01-02 15:04:05"),
		EndTime:   course.EndTime.Time.Format("2006-01-02 15:04:05"),
	})
	res.Teachers = make([]string, 0)
	for _, teacher := range course.Teachers {
		res.Teachers = append(res.Teachers, teacher.TeacherName)
	}
	response.OkWithData(res, c)
	return nil
}

type EnrollCourseRequest struct {
	CourseId int64 `json:"courseId"`
}

// EnrollCourse
// @Summary 报名课程接口
// @Tags 学生部分
// @Accept json
// @Param req body EnrollCourseRequest true "报名课程请求参数"
// @Success 200 {object} response.Response{data=GetCoursesByIdResponse} "success"
// @Router /api/user/courses/ [post]
func (u *UserApi) EnrollCourse(c *gin.Context) error {
	var req EnrollCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return err
	}
	if req.CourseId == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	baseClaims, exist := c.Get("claims")
	if !exist {
		fmt.Println("claims not exist")
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	claims := baseClaims.(*utils.CustomClaims)

	hasEnrolled, err := services.UserService.CountHasEnrolledCourses(claims.StudentPrimaryId, req.CourseId)
	if err != nil {
		return err
	}
	if hasEnrolled > 0 {
		return xerr.NewErrCodeMsg(3, "已经报名该课程")
	}
	err = services.UserService.EnrollCourse(claims.StudentPrimaryId, req.CourseId)
	if err != nil {
		return err
	}
	response.Ok(c)
	return nil
}

// DropCourse
// @Summary 退课接口
// @Tags 学生部分
// @Accept json
// @Param courseId path int true "课程id"
// @Success 200 {object} response.Response{data=UserGetCourseRequest} "success"
// @Router /api/user/courseId/{courseId} [delete]
func (u *UserApi) DropCourse(c *gin.Context) error {
	courseIdParam := c.Param("courseId")
	courseId, err := strconv.ParseInt(courseIdParam, 10, 64)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if courseId == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	baseClaims, exist := c.Get("claims")
	if !exist {
		fmt.Println("claims not exist")
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	claims := baseClaims.(*utils.CustomClaims)

	err = services.UserService.DropCourse(claims.StudentPrimaryId, courseId)
	response.Ok(c)
	return nil
}

type GetEnrolledCoursesRequest struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

type UserGetCoursesRequest struct {
	CourseId   int64    `json:"courseId"`
	CourseName string   `json:"courseName"`
	Capacity   int32    `json:"capacity"`
	Teachers   []string `json:"teachers"`
	Time       []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"time"`
	Location string `json:"location"`
}

// GetEnrolledCourses
// @Summary 查询已选课程信息
// @Tags 学生部分
// @Accept json
// @Success 200 {object} response.Response{data=UserGetCoursesRequest} "success"
// @Router /api/user/courses-selected [get]
func (u *UserApi) GetEnrolledCourses(c *gin.Context) error {
	baseClaims, exist := c.Get("claims")
	if !exist {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	claims := baseClaims.(*utils.CustomClaims)

	var student models.StudentModel
	if err := models.DB.Model(&models.StudentModel{}).
		Preload("Courses").
		Where("id = ?", claims.StudentPrimaryId).
		Find(&student).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	courseIds := make([]int64, 0)
	for _, course := range student.Courses {
		courseIds = append(courseIds, course.Id)
	}

	var courses []models.CourseModel
	if err := models.DB.Model(&models.CourseModel{}).
		Preload("Teachers").
		Where("id in (?)", courseIds).
		Find(&courses).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	var res []UserGetCoursesRequest
	teacherMap := make(map[int64][]string)
	for _, course := range courses {
		for _, teacher := range course.Teachers {
			teacherMap[course.Id] = append(teacherMap[course.Id], teacher.TeacherName)
		}
	}
	for _, course := range courses {
		tmp := UserGetCoursesRequest{
			CourseId:   course.Id,
			CourseName: course.CourseName,
			Capacity:   course.Capacity,
			Teachers:   teacherMap[course.Id],
			Time: make([]struct {
				StartTime string `json:"startTime"`
				EndTime   string `json:"endTime"`
			}, 0),
			Location: course.Location,
		}
		tmp.Time = append(tmp.Time, struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}{
			StartTime: course.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:   course.EndTime.Time.Format("2006-01-02 15:04:05"),
		})
		res = append(res, tmp)
	}
	response.OkWithData(res, c)
	return nil
}

// Logout
// @Summary 退出登录接口
// @Tags 学生部分
// @Accept json
// @Success 200 {object} response.Response{} "success"
// @Router /api/user/ [delete]
func (u *UserApi) Logout(c *gin.Context) error {
	response.Ok(c)
	return nil
}

type TestValidateRequest struct {
	Name     string `json:"name" binding:"required,timing"`
	Password string `json:"password" binding:"required,max=16,min=6"`
}

func (u *UserApi) TestValidate(c *gin.Context) error {
	var req TestValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			response.Fail(c)
			return nil
		}
		response.FailWithMessage(utils.ParseErr(errs), c)
		return nil
	}
	response.Ok(c)
	return nil
}
