package redisUtil

import (
	"fmt"
	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

type RedisInfo struct {
	Ip       string
	Port     int
	Password string
	Db       int
}

func RedisInit(redisInfo RedisInfo) {
	addr := fmt.Sprintf("%s:%v", redisInfo.Ip, redisInfo.Port)
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisInfo.Password,
		DB:       redisInfo.Db,
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic("redis ping error")
	}
}
