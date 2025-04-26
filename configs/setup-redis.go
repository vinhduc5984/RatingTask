package configs

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)


func ConnectRedis() *redis.Client{
	// connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // địa chỉ Redis server
		Password: "",               // mật khẩu (nếu có)
		DB:       0,                // sử dụng DB mặc định
	})
	// Kiểm tra kết nối
	ctx := context.Background()
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("aaaaa: ",pong) // Output: PONG nếu kết nối thành công
	return rdb
}


// Client instance
var RDB *redis.Client = ConnectRedis()