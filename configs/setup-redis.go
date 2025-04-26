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
	fmt.Println("aaaaa: ",pong) 
	return rdb
}


// Client instance
var RDB *redis.Client = ConnectRedis()


// handle some utils func for redis
func DeleteRDBByKey (ctx context.Context, keys []string){
	if len(keys) > 0{
		err := RDB.Del(ctx, keys...)
		if err != nil{
			fmt.Println("Delete redis database error: ",err.Err().Error())
		}
	}else{
		fmt.Println("Key is empty!")
	}
}