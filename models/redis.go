package models

import (
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func init() {
	redisURL := os.Getenv("REDIS_URL")
	password := ""
	if !strings.Contains(redisURL, "localhost") {
		parsedURL, _ := url.Parse(redisURL)
		password, _ = parsedURL.User.Password()
		redisURL = parsedURL.Host
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: password,
	})

	err := RedisClient.Ping().Err()
	if err != nil {
		log.Fatal(err)
	}
}
