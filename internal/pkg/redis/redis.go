package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type RedisClient struct {
	client *redis.Client
}

func New() RedisClient {
	c := RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
	return c
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) ExampleClient() error {
	ctx := context.Background()
	rdb := c.client

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		return err
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		return err
	}
	zap.L().Info("key", zap.String("value", val))

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		zap.L().Info("key2", zap.String("value", val))
		return nil
	}
	if err != nil {
		return err
	}
	fmt.Println("key2", val2)
	return nil
	// Output: key value
	// key2 does not exist
}
