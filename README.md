# Applied Rate Limiting in Go

This repository contains implementations of various rate-limiting algorithms in Go, integrated with Redis for distributed environments. Each algorithm is accompanied by visual illustrations to help understand its core logic and behavior.

## 🚀 Implemented Algorithms

### 1. Fixed Window Counter
The simplest form of rate limiting. It divides time into fixed-size windows (e.g., 60 seconds) and counts requests within each window.

- **Implementation**: [middleware/fixedWindow.middleware.go](middleware/fixedWindow.middleware.go)
- **Illustration**: [ilustrations/rate-limits/fixedWindow.excalidraw](ilustrations/rate-limits/fixedWindow.excalidraw)
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

## 🔍 Testing the Rate Limit

You can test the **Fixed Window** rate limit by hitting the `/get-details` endpoint. The current configuration allows **5 requests per 60 seconds**.

Using `curl`:
```bash
for i in {1..10}; do curl -i http://localhost:8080/get-details; echo; done
```

You should see `200 OK` for the first 5 requests and `429 Too Many Requests` for subsequent requests until the window resets.

## 📁 Project Structure

```text
.
├── config/             # Configuration (Redis client)
├── handlers/           # HTTP handlers
├── ilustrations/       # Excalidraw diagrams for algorithms
├── middleware/         # Rate limiting implementations
├── routes/             # Route definitions
├── docker-compose.yaml # Redis setup
└── main.go             # Entry point
```

## 🎨 Illustrations

The illustrations are created using [Excalidraw](https://excalidraw.com/). You can open the `.excalidraw` files directly in the Excalidraw web editor or use the VS Code Excalidraw extension to view them.
