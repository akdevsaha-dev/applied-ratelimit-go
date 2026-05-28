package middleware

import (
	"net/http"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
)

const (
	limit      = 5
	windowSize = 60 * time.Second
)

func FixedWindowRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var identifier string

		userId := r.Context().Value("userId")
		if userId != nil {
			identifier = "user:" + userId.(string)
		} else {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				ip = r.RemoteAddr
			}
			identifier = "ip:" + ip
		}

		key := "rate_limit:" + identifier

		count, err := config.RedisClient.Incr(config.Ctx, key).Result()

		if err != nil {
			http.Error(w, "Redis Error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			config.RedisClient.Expire(config.Ctx, key, windowSize)
		}

		if count > limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
