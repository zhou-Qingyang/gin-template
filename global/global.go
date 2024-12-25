package global

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	LocalCache = cache.New(5*time.Minute, 10*time.Minute)
)
