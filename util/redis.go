package util

import (
	"github.com/go-redis/redis/v8"
)

func Redis() *redis.Client {
	config,_ := GetConfig()
	host     := config["REDIS_HOST"]
	psw		 := config["REDIS_PASSWORD"]
	port 	 := config["REDIS_PORT"]

	rdb := redis.NewClient(&redis.Options{
		Addr:     host+":"+port,
		Password: psw, // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
