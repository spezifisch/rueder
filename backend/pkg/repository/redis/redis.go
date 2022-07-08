package redis

import (
	"context"

	"github.com/apex/log"
	"github.com/go-redis/redis/v8"
)

// RedisRepository internal state
type RedisRepository struct {
	rdb *redis.Client
	ctx context.Context
}

// NewRedisRepository returns a RedisRepository that wraps a redis DB
func NewRedisRepository(addr string, db int) *RedisRepository {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.WithError(err).WithField("addr", addr).Error("couldn't connect to redis")
		return nil
	}

	log.Info("connected to redis")

	return &RedisRepository{
		rdb: rdb,
		ctx: ctx,
	}
}
