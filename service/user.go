package service

import (
	"tz-gin/models"
	"tz-gin/models/common/response"
	"tz-gin/models/common/xerr"
	"tz-gin/utils"
)

type UserService struct {
}

func (u *UserService) FindByStudentId(studentId string) (res *models.UserModel, err error) {
	if err := models.DB.Model(&models.UserModel{}).
		Where("student_id =?", studentId).
		First(&res).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return res, nil
}

func (u *UserService) CreateAccount(userInfo *models.UserModel) (err error) {
	tx := models.DB.Begin()
	user := &models.UserModel{
		StudentId:   userInfo.StudentId,
		StudentName: userInfo.StudentName,
		Password:    utils.Md5(userInfo.Password),
		IsAdmin:     0, //默认不是管理员
	}
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Create(&models.StudentModel{
		UserId:      user.Id,
		StudentId:   userInfo.StudentId,
		StudentName: userInfo.StudentName,
		Courses:     nil,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (u *UserService) FindUserByUserId(userId int64) (res *models.UserModel, err error) {
	if err := models.DB.Model(&models.UserModel{}).
		Where("id =?", userId).
		First(&res).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return res, nil
}

func (u *UserService) ListCoursesBy(courseName string, location string, page int, limit int) (count int64, res []models.CourseModel, err error) {
	var size int64
	offset, limit, err := utils.GetPagination(page, limit)
	if err != nil {
		return 0, nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	coursesDb := models.DB.
		Model(&models.CourseModel{})
	if courseName != "" {
		coursesDb = coursesDb.Where("course_name LIKE ?", "%"+courseName+"%")
	}
	if location != "" {
		coursesDb = coursesDb.Where("location = ?", location)
	}
	var tmp []models.CourseModel
	if err := coursesDb.
		Count(&size).
		Offset(offset).
		Limit(limit).
		Find(&tmp).Error; err != nil {
		return 0, nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return size, tmp, nil
}

func (u *UserService) FindCourseByCourseId(courseId int) (res *models.CourseModel, err error) {
	if err := models.DB.Model(&models.CourseModel{}).
		Where("id = ?", courseId).
		Preload("Teachers").
		First(&res).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return res, nil
}

func (u *UserService) CountHasEnrolledCourses(userId int64, courseId int64) (count int64, err error) {
	if err := models.DB.
		Table("student_courses").
		Where("student_id = ? and course_id = ?", userId, courseId).
		Count(&count).Error; err != nil {
		return 0, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return count, nil
}

func (u *UserService) EnrollCourse(studentId int64, courseId int64) (err error) {
	if err := models.DB.Create(&models.StudentCourseModel{
		StudentId: studentId,
		CourseId:  courseId,
	}).Error; err != nil {
		return xerr.NewErrCodeMsg(3, "抢课失败")
	}
	return nil
}

func (u *UserService) DropCourse(studentId int64, courseId int64) (err error) {
	if err := models.DB.
		Table("student_courses").
		Where("student_id = ? and course_id = ?", studentId, courseId).
		Delete(&models.StudentCourseModel{}).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	return nil
}
