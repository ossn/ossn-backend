package models

import (
	"log"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	redisURL := strings.Split(os.Getenv("REDIS_URL"), "redis://")
	RedisClient = redis.NewClient(&redis.Options{Addr: redisURL[len(redisURL)-1]})
	err := RedisClient.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}
}
