package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type TeacherModel struct {
	Id          int64          `gorm:"primary_key" json:"id"`                                                            // 主键ID                                                     // ID
	TeacherName string         `json:"teacherName" form:"teacherName" gorm:"column:teacher_name;comment:教师名称;size:200;"` // 教师名称
	CreatedAt   sql.NullTime   `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`                                 // 创建时间
	UpdatedAt   sql.NullTime   `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
	Courses     []*CourseModel `json:"courses" gorm:"many2many:teacher_courses;foreignKey:Id;joinForeignKey:TeacherId;References:Id;joinReferences:CourseId"`
}

func (TeacherModel) TableName() string {
	return "teachers"
}
