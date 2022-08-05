package redisUtil

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var client *redis.Client

type RedisInfo struct {
	Ip       string
	Port     int
	Password string
	Db       int
}

func RedisInit(redisInfo RedisInfo) {
	addr := fmt.Sprintf("%s:%v", redisInfo.Ip, redisInfo.Port)
	client = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     redisInfo.Password,
		DB:           redisInfo.Db,
		PoolSize:     40,              // Redis连接池大小
		MaxRetries:   3,               // 最大重试次数
		IdleTimeout:  5 * time.Second, // 空闲链接超时时间
		MinIdleConns: 5,               // 空闲连接数量
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic("redis ping error")
	}
}

type Redis struct{}

func GetConnection() *redis.Client {
	return client
}

func (r Redis) Get(db int, key string) (string, error) {
	connection := GetConnection()
	connection.Do("select", db)
	return connection.Get(key).Result()
}
func (r Redis) Set(db int, key string, value interface{}, expiration time.Duration) error {
	connection := GetConnection()
	connection.Do("select", db)
	return connection.Set(key, value, expiration).Err()
}

func (r Redis) BatchSet(db int, kv map[string]interface{}, expiration time.Duration) error {
	connection := GetConnection()
	pipe := connection.Pipeline()
	pipe.Select(db)
	for k, v := range kv {
		pipe.Set(k, v, expiration)
	}
	_, err := pipe.Exec()
	return err
}
