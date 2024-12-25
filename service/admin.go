package service

import (
	"tz-gin/models"
	"tz-gin/models/common/response"
	"tz-gin/models/common/xerr"
	"tz-gin/utils"
)

type AdminService struct {
}

func (s *AdminService) ListTeachersByNames(names []string) (hasTeachers []models.TeacherModel, err error) {
	if err := models.DB.Model(&models.TeacherModel{}).
		Where("teacher_name IN (?)", names).
		Find(&hasTeachers).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return hasTeachers, nil
}

func (s *AdminService) DeleteCourseWithTeachers(courseId int64, hasTeachers []models.TeacherModel) (err error) {
	if err := models.DB.Model(&hasTeachers).
		Where("course_id = ?", courseId).
		Association("Courses").
		Clear(); err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	return nil
}

func (s *AdminService) ListCoursesByNames(courseName string, location string, page int, limit int) (hasCourses []models.CourseModel, count int64, err error) {
	offset, limit, err := utils.GetPagination(page, limit)
	if err != nil {
		return nil, 0, xerr.NewErrCode(response.SERVER_ERROR)
	}
	coursesDb := models.DB.
		Model(&models.CourseModel{})
	if courseName != "" {
		coursesDb = coursesDb.Where("course_name LIKE ?", "%"+courseName+"%")
	}
	if location != "" {
		coursesDb = coursesDb.Where("location = ?", location)
	}

	if err := coursesDb.
		Count(&count).
		Offset(offset).
		Limit(limit).
		Find(&hasCourses).Error; err != nil {
		return nil, 0, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return hasCourses, count, nil
}

func (s *AdminService) UpdateCourse(data models.CourseModel) (err error) {
	if err := models.DB.Begin().Model(&models.CourseModel{}).
		Where("id = ?", data.Id).
		Updates(map[string]interface{}{
			"course_name": data.CourseName,
			"capacity":    data.Capacity,
			"location":    data.Location,
			"start_time":  data.StartTime,
			"end_time":    data.EndTime,
		}).Error; err != nil {
		return xerr.NewErrCode(response.SERVER_ERROR)
	}
	return nil
}

// Page        int    `form:"page"`
// Limit       int    `form:"limit"`
// StudentName string `form:"studentName"`
// StudentId   string `form:"studentId"`
func (s *AdminService) ListStudents(page int, size int, studentName string, studentId string) (students []models.StudentModel, err error) {
	offset, limit, err := utils.GetPagination(page, size)
	if err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}

	db := models.DB.Model(&models.StudentModel{})
	if studentName != "" {
		db = db.Where("student_name LIKE ?", "%"+studentName+"%")
	}
	if studentId != "" {
		db = db.Where("student_id LIKE ?", "%"+studentId+"%")
	}
	if err := db.
		Offset(offset).
		Limit(limit).
		Find(&students).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return students, nil
}

func (s *AdminService) GetStudentById(id int64) (student *models.StudentModel, err error) {
	if err := models.DB.
		Model(&models.StudentModel{}).
		Where("id = ?", id).
		Preload("Courses").
		First(&student).Error; err != nil {
		return nil, xerr.NewErrCode(response.SERVER_ERROR)
	}
	return student, nil
}
