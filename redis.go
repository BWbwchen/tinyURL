package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

var rdb *redis.Client

func testConnectAndAutoConnect() {
	if err := rdb.Ping().Err(); err != nil {
		rdb = redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}
}

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	testConnectAndAutoConnect()
}

func RedisAdd(LongURL string, shortName string) {
	testConnectAndAutoConnect()
	err := rdb.Set(shortName, LongURL, 0).Err()
	if err != nil {
		fmt.Println("Redis Insert Error")
		panic(err)
	}
}

func RedisGet(shortName string) (string, error) {
	testConnectAndAutoConnect()
	return rdb.Get(shortName).Result()
}
