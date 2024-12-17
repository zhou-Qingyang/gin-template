package global

import (
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"time"
)

var (
	DBClient   *gorm.DB
	LocalCache = cache.New(5*time.Minute, 10*time.Minute)
)
