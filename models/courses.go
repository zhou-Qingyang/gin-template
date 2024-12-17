package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type CourseModel struct {
	Id         int64           `gorm:"primary_key" json:"id"`                                                         // 主键ID                                                     // ID
	CourseName string          `json:"courseName" form:"courseName" gorm:"column:course_name;comment:课程名称;size:200;"` // 课程名称
	StartTime  sql.NullTime    `json:"startTime" form:"startTime" gorm:"column:start_time;comment:开始时间;"`             // 开始时间
	EndTime    sql.NullTime    `json:"endTime" form:"endTime" gorm:"column:end_time;comment:结束时间;"`                   // 结束时间
	Location   string          `json:"location" form:"location" gorm:"column:location;comment:上课地点;size:200;"`        // 上课地点
	Capacity   int32           `json:"capacity" form:"capacity" gorm:"column:capacity;comment:课程容量;"`                 // 课程容量
	CreatedAt  sql.NullTime    `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`                              // 创建时间
	UpdatedAt  sql.NullTime    `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt  gorm.DeletedAt  `gorm:"index" json:"-"` // 删除时间
	Teachers   []*TeacherModel `json:"teachers" gorm:"many2many:teacher_courses;foreignKey:Id;joinForeignKey:CourseId;References:Id;joinReferences:TeacherId;"`
	Students   []*StudentModel `json:"students" gorm:"many2many:student_courses;foreignKey:Id;joinForeignKey:CourseId;References:Id;joinReferences:StudentId;"`
}

func (CourseModel) TableName() string {
	return "courses"
}
