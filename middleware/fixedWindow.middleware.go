package middleware

import (
	"net"
	"net/http"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
	"github.com/redis/go-redis/v9"
)

const (
	limit      = 5
	windowSize = 60 * time.Second
)

func FixedWindowRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var identifier string
		ctx := r.Context()
		userId := r.Context().Value("userId")
		if userId != nil {
			identifier = "user:" + userId.(string)
		} else {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil {
					ip = host
				} else {
					ip = r.RemoteAddr
				}
			}
			identifier = "ip:" + ip
		}

		key := "rate_limit:" + identifier

		pipe := config.RedisClient.Pipeline()

		incrCmd := pipe.Incr(ctx, key)
		ttlCmd := pipe.TTL(ctx, key)

		_, err := pipe.Exec(ctx)
		if err != nil && err != redis.Nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		count := incrCmd.Val()
		ttl := ttlCmd.Val()

		if ttl < 0 {
			config.RedisClient.Expire(ctx, key, windowSize)
		}

		if count > limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
