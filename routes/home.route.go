package routes

import (
	"net/http"

	"github.com/akdevsaha-dev/applied-ratelimit-go/handlers"
	"github.com/akdevsaha-dev/applied-ratelimit-go/middleware"
)

func RegisterHomeRoute(mux *http.ServeMux) {
	homeHandler := middleware.TokenBucketMiddleware(http.HandlerFunc(handlers.HomeHandler))
	mux.Handle("/get-details", homeHandler)
}
