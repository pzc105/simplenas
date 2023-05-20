package db

import (
	"fmt"
	"pnas/setting"

	"github.com/redis/go-redis/v9"
)

var GREDIS *redis.Client

func initRedis() {
	GREDIS = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", setting.GS.Redis.Ip, setting.GS.Redis.Port),
		Password: setting.GS.Redis.Password,
		DB:       0,
	})
}
