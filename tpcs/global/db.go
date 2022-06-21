package global

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

var (
	DBEngine    *gorm.DB
	RedisClient *redis.Client
)
