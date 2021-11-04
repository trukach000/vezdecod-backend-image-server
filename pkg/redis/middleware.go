package redis

import (
	"net/http"

	"github.com/go-redis/redis"
)

func NewRedisMiddleware(r *redis.Client, ignoreRedis bool) *RedisMiddleware {
	return &RedisMiddleware{
		client: r,
		ignore: ignoreRedis,
	}
}

type RedisMiddleware struct {
	client *redis.Client
	ignore bool
}

func (rm *RedisMiddleware) Attach(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = PutRedisToContext(rm.client, ctx)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
