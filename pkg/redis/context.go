package redis

import (
	"context"
	"errors"

	"github.com/go-redis/redis"
)

const (
	missingRedisInstanceError = "missing redis instance"
	REDIS_INSTANCE_CTX_NAME   = "redisInstance"
)

func GetRedisFromContext(ctx context.Context) (*redis.Client, error) {
	dbEntry, ok := ctx.Value(REDIS_INSTANCE_CTX_NAME).(*redis.Client)
	if !ok {
		return nil, errors.New(REDIS_INSTANCE_CTX_NAME)
	}
	return dbEntry, nil
}

func PutRedisToContext(entry *redis.Client, ctx context.Context) context.Context {
	return context.WithValue(ctx, REDIS_INSTANCE_CTX_NAME, entry)
}
