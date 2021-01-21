package util

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type MyRedis struct {
	Host string
	Password string
	DB int
	Port string
	Client *redis.Client
}

func NewRedis() *MyRedis {
	config,_ := GetConfig()

	redisClient := &MyRedis{
		Host:     config["REDIS_HOST"],
		Password: config["REDIS_PASSWORD"],
		DB:       0,
		Port:	  config["REDIS_PORT"],
		Client:   nil,
	}

	redisClient.Client = redis.NewClient(&redis.Options{
		Addr:     redisClient.Host+":"+redisClient.Port,
		Password: redisClient.Password, // no password set
		DB:       0,  // use default DB
	})

	return redisClient
}

/**
 * 监听队列
 */
func (redis *MyRedis) OnQueue(wg *sync.WaitGroup,queue string,timeOut time.Duration,call func(*redis.StringSliceCmd,error))  {
	defer wg.Done()
	for {
		var ctx = context.Background()
		result := redis.Client.BLPop(ctx,timeOut,queue)
		call(result,nil)
	}
}