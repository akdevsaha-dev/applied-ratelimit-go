# Applied Rate Limiting in Go

This repository contains implementations of various rate-limiting algorithms in Go, integrated with Redis for distributed environments. Each algorithm is designed to handle different traffic patterns and use cases.

## 🚀 Implemented Algorithms

### 1. Fixed Window Counter
The simplest form of rate limiting. It divides time into fixed-size windows (e.g., 60 seconds) and counts requests within each window.
- **Config**: 5 requests per 60 seconds.
- **Implementation**: [middleware/fixedWindow.middleware.go](middleware/fixedWindow.middleware.go)
- **Status**: ✅ Implemented

### 2. Token Bucket
Tokens are added to a bucket at a fixed rate. Each request consumes a token. If the bucket is empty, the request is rejected. This allows for bursts of traffic up to the bucket capacity.
- **Config**: Capacity of 10 tokens, refill rate of 1 token per second.
- **Implementation**: [middleware/tokenBucket.middleware.go](middleware/tokenBucket.middleware.go)
- **Status**: ✅ Implemented

### 3. Leaky Bucket
Requests enter a bucket and "leak" out at a constant rate. This smoothes out bursts of traffic and ensures a steady processing rate.
- **Config**: Capacity of 10, leak rate of 1 per second.
- **Implementation**: [middleware/leakyBucket.middleware.go](middleware/leakyBucket.middleware.go)
- **Status**: ✅ Implemented

### 4. Sliding Window Log
Tracks every request's timestamp in a sorted set. To check the rate limit, it removes old timestamps outside the current window and counts the remaining ones. Highly accurate but memory-intensive.
- **Config**: 5 requests per 60 seconds.
- **Implementation**: [middleware/slidingLog.middleware.go](middleware/slidingLog.middleware.go)
- **Status**: ✅ Implemented

### 5. Sliding Window Counter
A hybrid approach that uses the current window's count and a weighted count from the previous window to estimate the current request rate without the memory overhead of the Sliding Window Log.
- **Config**: 5 requests per 60 seconds.
- **Implementation**: [middleware/slidingCounter.middleware.go](middleware/slidingCounter.middleware.go)
- **Status**: ✅ Implemented

## 🛠️ Prerequisites

- [Go](https://go.dev/doc/install) (1.21+)
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Redis](https://redis.io/) (via Docker)

## 🏃 Getting Started

1. **Clone the repository**:
   ```bash
   git clone https://github.com/akdevsaha-dev/applied-ratelimit-go.git
   cd applied-ratelimit-go
   ```

2. **Start the Redis container**:
   ```bash
   docker-compose up -d
   ```

3. **Run the Go application**:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## 🔄 Switching Algorithms

To test a different rate-limiting algorithm, update the middleware in `routes/home.route.go`:

```go
// routes/home.route.go

func RegisterHomeRoute(mux *http.ServeMux) {
    // Switch between:
    // middleware.FixedWindowRateLimit
    // middleware.TokenBucketMiddleware
    // middleware.LeakyBucketMiddleware
    // middleware.SlidingLogMiddleware
    // middleware.SlidingCounterMiddleware
    
    homeHandler := middleware.TokenBucketMiddleware(http.HandlerFunc(handlers.HomeHandler))
    mux.Handle("/get-details", homeHandler)
}
```

## 🔍 Testing the Rate Limit

You can test the active rate limit by hitting the `/get-details` endpoint.

Using `curl`:
```bash
for i in {1..12}; do curl -i http://localhost:8080/get-details; echo; sleep 0.5; done
```

Depending on the algorithm, you will see `200 OK` for allowed requests and `429 Too Many Requests` when the limit is reached.

## 📁 Project Structure

```text
.
├── config/             # Redis client configuration
├── handlers/           # HTTP handlers
├── illustrations/      # Excalidraw diagrams for algorithms
├── middleware/         # Rate limiting implementations
├── routes/             # Route definitions
├── docker-compose.yaml # Redis setup
└── main.go             # Application entry point
```

## 🎨 Illustrations

The illustrations are created using [Excalidraw](https://excalidraw.com/). You can find them in the `illustrations/rate-limits/` directory. Open them directly in the Excalidraw web editor or use the VS Code Excalidraw extension.
