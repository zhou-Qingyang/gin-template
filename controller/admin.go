package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"strconv"
	"tz-gin/models"
	"tz-gin/models/common/response"
	"tz-gin/models/common/xerr"
	"tz-gin/utils"
)

type AdminApi struct{}

type AddCourseRequest struct {
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

// AddCourse
// @Summary 添加课程接口
// @Description 添加课程接口
// @Tags 管理员部分
// @Accept json
// @Param req body AddCourseRequest true "学生注册请求参数"
// @Success 200 {object} response.Response{data=map[string]interface{}} "success"
// @Router /api/admin/courses [post]
func (a *AdminApi) AddCourse(c *gin.Context) error {
	var req AddCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if len(req.CourseName) == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if len(req.Teachers) == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	startTime, endTime := sql.NullTime{Valid: true}, sql.NullTime{Valid: true}
	forMatStartTime, err := utils.ParseDate(req.Time[0].StartTime)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	forMatEndTime, err := utils.ParseDate(req.Time[0].EndTime)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	if forMatStartTime.After(forMatEndTime) {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	startTime.Time = forMatStartTime
	endTime.Time = forMatEndTime

	hasTeachers, err := services.AdminService.ListTeachersByNames(req.Teachers)
	if err != nil {
		return err
	}

	if len(hasTeachers) != len(req.Teachers) {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	// 先创建一个课程
	course := models.CourseModel{
		CourseName: req.CourseName,
		Capacity:   req.Capacity,
		Location:   req.Location,
		StartTime:  startTime,
		EndTime:    endTime,
	}

	tx := models.DB.Begin()
	if err := tx.Create(&course).Error; err != nil {
		tx.Rollback()
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	// 遍历每个老师并建立关联关系
	for _, teacher := range hasTeachers {
		if err := tx.Model(&teacher).
			Association("Courses").Append(&course); err != nil {
			tx.Rollback()
			return xerr.NewErrCode(response.SERVER_ERROR)
		}
	}
	tx.Commit()

	response.OkWithData(map[string]interface{}{
		"id": course.Id,
	}, c)
	return nil
}

// DeleteCourse
// @Summary 删除课程接口
// @Description 删除课程接口
// @Tags 管理员部分
// @Accept json
// @Param courseId path int true "课程id"
// @Success 200 {object} response.Response{} "success"
// @Router /api/admin/courses/{courseId} [delete]
func (a *AdminApi) DeleteCourse(c *gin.Context) error {
	courseIdStr := c.Param("courseId")
	var teachers []models.TeacherModel
	if err := models.DB.Model(&models.TeacherModel{}).
		Find(&teachers).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	courseId, err := strconv.ParseInt(courseIdStr, 10, 64)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	err = services.AdminService.DeleteCourseWithTeachers(courseId, teachers)
	if err != nil {
		return err
	}
	response.Ok(c)
	return nil
}

// UpdateCourse
// @Summary 修改课程接口
// @Description 修改课程接口
// @Tags 管理员部分
// @Accept json
// @Param req body AddCourseRequest true "修改课程请求参数"
// @Success 200 {object} response.Response{} "success"
// @Router /api/admin/courses [put]
func (a *AdminApi) UpdateCourse(c *gin.Context) error {
	var req AddCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if len(req.CourseName) == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if len(req.Teachers) == 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	if req.CourseId <= 0 {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	startTime, endTime := sql.NullTime{Valid: true}, sql.NullTime{Valid: true}
	forMatStartTime, err := utils.ParseDate(req.Time[0].StartTime)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	forMatEndTime, err := utils.ParseDate(req.Time[0].EndTime)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	if forMatStartTime.After(forMatEndTime) {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	startTime.Time = forMatStartTime
	endTime.Time = forMatEndTime

	hasTeachers, err := services.AdminService.ListTeachersByNames(req.Teachers)
	if err != nil {
		return err
	}

	var toAssignTeachers []models.TeacherModel
	for _, te := range hasTeachers {
		if utils.HasContainInSlice(te.TeacherName, req.Teachers) {
			toAssignTeachers = append(toAssignTeachers, te)
		}
	}

	if len(toAssignTeachers) != len(req.Teachers) {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}

	tx := models.DB.Begin()
	// 删除满足条件的所有关联关系
	if err := tx.Model(&hasTeachers).
		Where("course_id = ?", req.CourseId).
		Association("Courses").
		Clear(); err != nil {
		tx.Rollback()
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	if err := tx.Delete(&models.CourseModel{}, req.CourseId).Error; err != nil {
		tx.Rollback()
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	err = services.AdminService.UpdateCourse(models.CourseModel{
		CourseName: req.CourseName,
		Capacity:   req.Capacity,
		Location:   req.Location,
		StartTime:  startTime,
		EndTime:    endTime,
		Id:         req.CourseId,
	})
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, teacher := range toAssignTeachers {
		if err := tx.Model(&teacher).
			Association("Courses").Append(&models.CourseModel{Id: req.CourseId}); err != nil {
			tx.Rollback()
			return xerr.NewErrCode(response.SERVER_ERROR)
		}
	}
	tx.Commit()
	response.Ok(c)
	return nil
}

type CoursesRequest struct {
	CourseName string   `form:"courseName"`
	Teachers   []string `form:"teachers"`
	Page       int      `form:"page"`
	Limit      int      `form:"limit"`
	Location   string   `form:"location"`
}

type CourseItem struct {
	Id         int64  `json:"id"`
	CourseName string `json:"courseName"`
	Capacity   int32  `json:"capacity"`
	Location   string `json:"location"`
	Time       []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"time"`
	Teachers []string `json:"teachers"`
}

type CoursesResponse struct {
	Rows []CourseItem `json:"rows"`
	Size int          `json:"size"`
}

// GetCourses
// @Summary 管理员获取课程列表
// @Description 管理员获取课程列表
// @Tags 管理员部分
// @Accept json
// @Param req query CoursesRequest true "修改课程请求参数"
// @Success 200 {object} response.Response{data=CoursesResponse} "success"
// @Router /api/admin/courses [get]
func (a *AdminApi) GetCourses(c *gin.Context) error {
	var req CoursesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return err
	}
	tmp, count, err := services.AdminService.ListCoursesByNames(req.CourseName, req.Location, req.Page, req.Limit)
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
		var time []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}
		time = append(time, struct {
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
			Time:       time,
			Teachers:   courseTeacherMap[course.Id],
		})
	}

	response.OkWithData(CoursesResponse{
		Size: int(count),
		Rows: rows,
	}, c)
	return nil
}

type GetCoursesByIdRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type GetCoursesByIdResponse struct {
	Id         int64  `json:"id"`
	CourseName string `json:"courseName"`
	Capacity   int32  `json:"capacity"`
	Location   string `json:"location"`
	Time       []struct {
		StartTime string `json:"startTime"`
		EndTime   string `json:"endTime"`
	} `json:"time"`
	Teachers []string `json:"teachers"`
	Students []struct {
		Name      string `json:"name"`
		StudentId string `json:"studentId"`
	} `json:"students"`
	TotalStudents int `json:"totalStudents"`
}

// GetCoursesById
// @Summary 管理员获取课程详情
// @Description 管理员获取课程详情
// @Tags 管理员部分
// @Accept json
// @Param courseId path int true "课程id"
// @Param req query GetCoursesByIdRequest true "修改课程请求参数"
// @Success 200 {object} response.Response{data=GetCoursesByIdResponse} "success"
// @Router /api/admin/courses/{courseId} [get]
func (a *AdminApi) GetCoursesById(c *gin.Context) error {
	courseIdParam := c.Param("courseId")
	courseId, err := strconv.Atoi(courseIdParam)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	var req GetCoursesByIdRequest
	if err := c.BindQuery(&req); err != nil {
		return err
	}

	var size int64
	offset, limit, err := utils.GetPagination(req.Page, req.Limit)
	if err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	course, err := services.UserService.FindCourseByCourseId(courseId)
	if err != nil {
		return err
	}

	if err = models.DB.
		Model(&models.StudentModel{}).
		Table("students as s").
		Joins("JOIN student_courses as sc ON sc.student_id = s.id").
		Joins("JOIN courses as c ON c.id = sc.course_id").
		Where("c.id = ?", courseId).
		Count(&size).
		Offset(offset).
		Limit(limit). // 限制只加载2个教师
		Find(&course.Students).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
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

	for _, student := range course.Students {
		res.Students = append(res.Students, struct {
			Name      string `json:"name"`
			StudentId string `json:"studentId"`
		}{
			Name:      student.StudentName,
			StudentId: student.StudentId,
		})
	}
	res.TotalStudents = int(size)
	response.OkWithData(res, c)
	return nil
}

type GetStudentsRequest struct {
	Page        int    `form:"page"`
	Limit       int    `form:"limit"`
	StudentName string `form:"studentName"`
	StudentId   string `form:"studentId"`
}

type GetStudentsResponse struct {
	Students []struct {
		StudentId    string `json:"studentId"`
		StudentName  string `json:"studentName"`
		TotalCourses int    `json:"totalCourses"`
	} `json:"students"`
}

// GetStudents
// @Summary 管理员获取学生列表
// @Description 管理员获取学生列表
// @Tags 管理员部分
// @Accept json
// @Param req query GetStudentsRequest true "修改课程请求参数"
// @Success 200 {object} response.Response{data=GetStudentsResponse} "success"
// @Router /api/admin/students [get]
func (a *AdminApi) GetStudents(c *gin.Context) error {
	var req GetStudentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return err
	}

	students, err := services.AdminService.ListStudents(req.Page, req.Limit, req.StudentName, req.StudentId)
	if err != nil {
		return err
	}
	var res []struct {
		StudentId    string `json:"studentId"`
		StudentName  string `json:"studentName"`
		TotalCourses int    `json:"totalCourses"`
	}

	studentsIds := make([]int64, 0)
	for _, student := range students {
		studentsIds = append(studentsIds, student.Id)
	}

	// 查询所有课程
	var courses []models.CourseModel
	if err := models.DB.
		Model(&models.CourseModel{}).
		Preload("Students", models.DB.Where("id IN (?)", studentsIds)).
		Find(&courses).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}

	var courseMap = make(map[int64]int)
	for _, course := range courses {
		for _, student := range course.Students {
			if utils.HasContainInSliceInt64(student.Id, studentsIds) {
				courseMap[student.Id] += 1
			}
		}
	}

	for _, student := range students {
		res = append(res, struct {
			StudentId    string `json:"studentId"`
			StudentName  string `json:"studentName"`
			TotalCourses int    `json:"totalCourses"`
		}{
			StudentId:    student.StudentId,
			StudentName:  student.StudentName,
			TotalCourses: courseMap[student.Id],
		})
	}

	response.OkWithData(GetStudentsResponse{Students: res}, c)
	return nil
}

type GetStudentResponse struct {
	StudentName string `json:"studentName"`
	Courses     []struct {
		Id         int64  `json:"id"`
		CourseName string `json:"courseName"`
		Capacity   int32  `json:"capacity"`
		Location   string `json:"location"`
		Time       []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		} `json:"time"`
		Teachers []string `json:"teachers"`
	}
}

// GetStudent
// @Summary 管理员获取学生详情
// @Description 管理员获取学生详情
// @Tags 管理员部分
// @Accept json
// @Param studentId path int true "学生id"
// @Success 200 {object} response.Response{data=GetStudentResponse} "success"
// @Router /api/admin/students/{studentId} [get]
func (a *AdminApi) GetStudent(c *gin.Context) error {
	studentIdParam := c.Param("studentId")
	studentIdStr, err := strconv.Atoi(studentIdParam)
	if err != nil {
		return xerr.NewErrCode(response.PARAMETER_ERROR)
	}
	student, err := services.AdminService.GetStudentById(int64(studentIdStr))
	if err != nil {
		return err
	}

	courseIds := make([]int64, 0)
	for _, course := range student.Courses {
		courseIds = append(courseIds, course.Id)
	}
	var studentCourses []models.CourseModel
	if err := models.DB.
		Model(&models.CourseModel{}).
		Preload("Teachers").
		Where("id IN (?)", courseIds).
		Find(&studentCourses).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	var teacherMap = make(map[int64][]string)
	for _, course := range studentCourses {
		for _, teacher := range course.Teachers {
			teacherMap[course.Id] = append(teacherMap[course.Id], teacher.TeacherName)
		}
	}
	var res GetStudentResponse
	res.StudentName = student.StudentName
	for _, course := range student.Courses {
		var time []struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}
		time = append(time, struct {
			StartTime string `json:"startTime"`
			EndTime   string `json:"endTime"`
		}{
			StartTime: course.StartTime.Time.Format("2006-01-02 15:04:05"),
			EndTime:   course.EndTime.Time.Format("2006-01-02 15:04:05"),
		})
		res.Courses = append(res.Courses, struct {
			Id         int64  `json:"id"`
			CourseName string `json:"courseName"`
			Capacity   int32  `json:"capacity"`
			Location   string `json:"location"`
			Time       []struct {
				StartTime string `json:"startTime"`
				EndTime   string `json:"endTime"`
			} `json:"time"`
			Teachers []string `json:"teachers"`
		}{
			Id:         course.Id,
			CourseName: course.CourseName,
			Capacity:   course.Capacity,
			Location:   course.Location,
			Time:       time,
			Teachers:   teacherMap[course.Id],
		})
	}

	response.OkWithData(res, c)
	return nil
}
