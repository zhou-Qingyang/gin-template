package global

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"tz-gin/models"
)

func GormMysql() *gorm.DB {
	// MySQL 配置
	mysqlConfig := mysql.Config{
		DSN:                       "root:123456@tcp(127.0.0.1:3306)/tenzor2024?charset=utf8mb4&parseTime=True&loc=Local", // DSN 数据源名称
		DefaultStringSize:         191,                                                                                   // string 类型字段的默认长度
		SkipInitializeWithVersion: false,                                                                                 // 根据版本自动配置
	}

	// 打开数据库连接
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // 设置 GORM 日志级别为 Info
	})
	if err != nil {
		log.Fatalf("MySQL connection failed: %v\n", err)
		return nil
	}

	// 设置数据库连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get raw DB object: %v\n", err)
		return nil
	}

	// 设置连接池最大打开连接数
	sqlDB.SetMaxOpenConns(50)
	// 设置连接池最大空闲连接数
	sqlDB.SetMaxIdleConns(20)
	// 设置连接最大可复用时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 输出成功日志
	fmt.Println("MySQL has initialized")
	//自动迁移（根据需要选择是否使用）
	err = db.AutoMigrate(&models.CourseModel{}, &models.StudentModel{}, &models.TeacherModel{}, &models.UserModel{}) // 可以进行自动迁移
	if err != nil {
		log.Fatalf("Failed to migrate: %v\n", err)
	}
	return db
}
