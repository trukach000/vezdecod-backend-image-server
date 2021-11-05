package redisclient

import (
	"net/http"

	"github.com/go-redis/redis"
)

func NewRedisMiddleware(r *redis.Client) *RedisMiddleware {
	return &RedisMiddleware{
		client: r,
	}
}

type RedisMiddleware struct {
	client *redis.Client
}

func (rm *RedisMiddleware) Attach(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = PutRedisToContext(rm.client, ctx)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
