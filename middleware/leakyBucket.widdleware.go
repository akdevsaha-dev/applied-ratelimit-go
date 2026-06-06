package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/akdevsaha-dev/applied-ratelimit-go/config"
)

type LeakyBucketState struct {
	Water     float64 `json:"water"`
	LastCheck int64   `json:"last_check"`
}

func LeakyBucketMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		const (
			capacity = 10.0
			leakRate = 1.0
		)

		ctx := context.Background()

		identifier := getIdentifier(r)

		key := "leaky_bucket:" + identifier

		now := time.Now().Unix()

		var state LeakyBucketState

		val, err := config.RedisClient.Get(ctx, key).Result()

		if err != nil {

			state = LeakyBucketState{
				Water:     0,
				LastCheck: now,
			}

		} else {

			_ = json.Unmarshal([]byte(val), &state)

			elapsed := float64(
				now - state.LastCheck,
			)

			state.Water -= elapsed * leakRate

			if state.Water < 0 {
				state.Water = 0
			}

			state.LastCheck = now
		}

		if state.Water+1 > capacity {

			http.Error(
				w,
				"Rate Limit Exceeded",
				http.StatusTooManyRequests,
			)

			return
		}

		state.Water++

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
