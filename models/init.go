package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
	"tz-gin/config"
	dblog "tz-gin/logger"
)

var DB *gorm.DB

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&&parseTime=True&loc=Local",
		config.Config.MysqlUser,
		config.Config.MysqlPass,
		config.Config.MysqlHost,
		config.Config.MysqlPort,
		config.Config.MysqlName)
	var dbLogger logger.Interface
	if dblog.DatabaseLogger == nil {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		logLevels := map[string]int{
			"error": 2,
			"warn":  3,
			"info":  4,
		}

		levels, ok := logLevels[config.Config.LogLevel]
		if !ok {
			levels = 4
		}
		dbLogger = logger.New(
			log.New(dblog.DataLogger{dblog.DatabaseLogger}, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             200 * time.Millisecond,
				LogLevel:                  logger.LogLevel(levels),
				IgnoreRecordNotFoundError: true,
				Colorful:                  false,
			},
		)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: dbLogger})
	if err != nil {
		panic(err)
	}

	DB = db

	if !config.Config.AppProd {
		db.AutoMigrate(&CourseModel{}, &StudentModel{}, &TeacherModel{}, &UserModel{}) // 可以进行自动迁移
	}
}
