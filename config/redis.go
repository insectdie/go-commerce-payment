package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConnection struct {
	Host string
	Port string
	Pass string
	DB   int
}

func ConnectToRedis(conn RedisConnection) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conn.Host, conn.Port),
		Password: conn.Pass,
		DB:       conn.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to redis: %w", err)
	}

	return client, nil
}
