package models

type StudentCourseModel struct {
	StudentId int64 `json:"studentId" form:"studentId" gorm:"column:student_id;comment:学生ID;size:200;"` // 学生ID
	CourseId  int64 `json:"courseId" form:"courseId" gorm:"column:course_id;comment:课程ID;size:200;"`    // 课程ID
}

func (StudentCourseModel) TableName() string {
	return "student_courses"
}
