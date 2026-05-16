package redis

import (
    "context"
    "log"

    "github.com/redis/go-redis/v9"
)
type RedisClient struct {
    *redis.Client
}
func NewRedisClient() RedisClient {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    err := client.Ping(context.Background()).Err()
    if err != nil {
        log.Fatalf("failed to connect redis: %v", err)
    }

    log.Println("redis connected successfully")

    return RedisClient{client}
}