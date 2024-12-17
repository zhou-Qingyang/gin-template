package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type UserModel struct {
	Id          int64          `gorm:"primary_key" json:"id"`
	StudentId   string         `json:"studentId" form:"studentId" gorm:"column:student_id;comment:学生ID;size:200;"`       // 学生ID
	StudentName string         `json:"studentName" form:"studentName" gorm:"column:student_name;comment:学生姓名;size:200;"` // 学生姓名                                                     // ID
	Password    string         `json:"password" form:"password" gorm:"column:password;comment:密码;size:200;"`             // 密码
	IsAdmin     int8           `json:"isAdmin" gorm:"column:is_admin;comment:是否管理员;"`                                    // 是否管理员
	CreatedAt   sql.NullTime   `json:"createdAt" gorm:"column:created_at;comment:创建时间;"`                                 // 创建时间
	UpdatedAt   sql.NullTime   `json:"updatedAt" gorm:"column:updated_at;comment:更新时间;"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // 删除时间
}

func (UserModel) TableName() string {
	return "users"
}
