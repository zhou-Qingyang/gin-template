package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type StudentModel struct {
	Id          int64          `gorm:"primary_key" json:"id"`
	StudentId   string         `json:"studentId" form:"studentId" gorm:"column:student_id;comment:学生ID;size:200;"`       // 学生ID
	StudentName string         `json:"studentName" form:"studentName" gorm:"column:student_name;comment:学生姓名;size:200;"` // 学生姓名
	UserId      int64          `json:"userId" form:"userId" gorm:"column:user_id;comment:用户ID;"`                         // 用户ID
	CreatedAt   sql.NullTime   `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`                                 // 创建时间
	UpdatedAt   sql.NullTime   `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
	Courses     []*CourseModel `json:"courses" gorm:"many2many:student_courses;foreignKey:Id;joinForeignKey:StudentId;References:Id;joinReferences:CourseId"`
}

func (StudentModel) TableName() string {
	return "students"
}
