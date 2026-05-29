package middleware

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func SlidingLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			limit      = 5
			windowSize = 60 * time.Second
		)

		ctx := r.Context()

		var identifier string

		userId := ctx.Value("userId")

		if userId != nil {
			identifier = "userId" + userId.(string)
		} else {
			ip := r.Header.Get("X-FORWARDED-FOR")
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

		key := "rate_limit" + identifier

		now := time.Now().UnixMilli()

		windowStart := now - windowSize.Milliseconds()

		pipe := config.RedisClient.Pipeline()

		pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprint(windowStart))

		countCmd := pipe.ZCard(ctx, key)

		_, err := pipe.Exec(ctx)

		if err != nil && err != redis.Nil {
			http.Error(
				w,
				"Internal server error",
				http.StatusInternalServerError,
			)
		}
		count := countCmd.Val()
		if count >= limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		requestId := uuid.NewString()

		pipe = config.RedisClient.TxPipeline()
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: requestId,
		})
		_, err = pipe.Exec(ctx)

		if err != nil && err != redis.Nil {
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
