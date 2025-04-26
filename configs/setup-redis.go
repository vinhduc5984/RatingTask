package configs

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

func ConnectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		// Production mode: connect via REDIS_URL
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("Failed to parse REDIS_URL: %v", err)
		}
		RedisClient = redis.NewClient(opt)
	} else {
		// Local mode: connect to localhost:6379
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password for local
			DB:       0,  // default DB
		})
	}

	// Check connection
	err := pingRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully!")
	return RedisClient
}

func pingRedis() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("Redis PING response:", pong)
	return nil
}

// Client instance
var RDB *redis.Client = ConnectRedis()
