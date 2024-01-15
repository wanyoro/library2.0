package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"fmt"

	
)


var (RedisClient *redis.Client)
//var Ctx = context.Background()

func StartRedis()(context.Context){
	RedisClient = redis.NewClient(&redis.Options{})
	fmt.Println("Redis is here")

	ctx:= context.Background()

	RedisClient.FlushAll(ctx)

	RedisClient.Set(ctx, "hello", "world",0)

	return ctx
}

func SET(key string, value string)error{
	ctx:= StartRedis()
	err:=RedisClient.Set(ctx, key, value, 0).Err()
	return err
}

func GET(key string) (string, error){
	ctx:= StartRedis()
	val,err := RedisClient.Get(ctx,key).Result()
	return val, err
}