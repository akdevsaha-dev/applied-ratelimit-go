package middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
)

func SlidingCounterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const limit = 5

		windowSize := 60 * time.Second
		windowSizeMs := windowSize.Milliseconds()

		ctx := r.Context()

		// Determine identifier (User ID or IP)
		var identifier string

		if userID := ctx.Value("userId"); userID != nil {
			identifier = "user:" + userID.(string)
		} else {
			ip := r.Header.Get("X-Forwarded-For")

			if ip != "" {
				// Use first IP if multiple proxies are involved
				ip = strings.TrimSpace(strings.Split(ip, ",")[0])
			} else {
				host, _, err := net.SplitHostPort(r.RemoteAddr)
				if err == nil {
					ip = host
				} else {
					ip = r.RemoteAddr
				}
			}

			identifier = "ip:" + ip
		}

		// Current time in milliseconds
		now := time.Now().UnixMilli()

		// Calculate current and previous windows
		currentWindow := now / windowSizeMs
		previousWindow := currentWindow - 1

		currentWindowStart := currentWindow * windowSizeMs

		// Time elapsed inside current window
		elapsed := now - currentWindowStart

		// Weight of previous window
		weight := float64(windowSizeMs-elapsed) / float64(windowSizeMs)

		currentKey := fmt.Sprintf(
			"rate_limit:%s:%d",
			identifier,
			currentWindow,
		)

		previousKey := fmt.Sprintf(
			"rate_limit:%s:%d",
			identifier,
			previousWindow,
		)

		// Fetch current and previous counts
		pipe := config.RedisClient.Pipeline()

		currentCmd := pipe.Get(ctx, currentKey)
		previousCmd := pipe.Get(ctx, previousKey)

		_, err := pipe.Exec(ctx)
		if err != nil && err.Error() != "redis: nil" {
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			return
		}

		currentCount := 0
		previousCount := 0

		if val, err := currentCmd.Int(); err == nil {
			currentCount = val
		}

		if val, err := previousCmd.Int(); err == nil {
			previousCount = val
		}

		// Sliding window formula
		totalRequests := float64(previousCount)*weight + float64(currentCount)

		if totalRequests >= float64(limit) {
			w.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
			w.Header().Set("X-RateLimit-Remaining", "0")

			http.Error(
				w,
				"Rate Limit Exceeded",
				http.StatusTooManyRequests,
			)
			return
		}

		// Increment current window counter
		txPipe := config.RedisClient.TxPipeline()

		txPipe.Incr(ctx, currentKey)
		txPipe.Expire(ctx, currentKey, windowSize*2)

		_, err = txPipe.Exec(ctx)
		if err != nil {
			http.Error(
				w,
				"Internal Server Error",
				http.StatusInternalServerError,
			)
			return
		}

		remaining := limit - int(totalRequests) - 1
		if remaining < 0 {
			remaining = 0
		}

		w.Header().Set(
			"X-RateLimit-Limit",
			strconv.Itoa(limit),
		)

		w.Header().Set(
			"X-RateLimit-Remaining",
			strconv.Itoa(remaining),
		)

		next.ServeHTTP(w, r)
	})
}