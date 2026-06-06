package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
)

type BucketState struct {
	Tokens     float64 `json:"tokens"`
	LastRefill int64   `json:"last_refill"`
}

func TokenBucketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const (
			capacity   = 10.0
			refillRate = 1.0
		)

		ctx := r.Context()

		identifier := getIdentifier(r)

		key := "token_bucket:" + identifier

		now := time.Now().Unix()

		var state BucketState

		val, err := config.RedisClient.Get(ctx, key).Result()

		if err != nil {

			state = BucketState{
				Tokens:     capacity,
				LastRefill: now,
			}

		} else {

			_ = json.Unmarshal([]byte(val), &state)

			elapsed := float64(now - state.LastRefill)

			state.Tokens += elapsed * refillRate

			if state.Tokens > capacity {
				state.Tokens = capacity
			}

			state.LastRefill = now
		}

		if state.Tokens < 1 {

			http.Error(
				w,
				"Rate Limit Exceeded",
				http.StatusTooManyRequests,
			)

			return
		}

		state.Tokens--

		data, _ := json.Marshal(state)

		config.RedisClient.Set(
			ctx,
			key,
			data,
			24*time.Hour,
		)

		next.ServeHTTP(w, r)
	})
}

func getIdentifier(r *http.Request) string {

	ip := r.Header.Get("X-Forwarded-For")

	if ip != "" {
		return strings.TrimSpace(
			strings.Split(ip, ",")[0],
		)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)

	if err == nil {
		return host
	}

	return r.RemoteAddr
}
