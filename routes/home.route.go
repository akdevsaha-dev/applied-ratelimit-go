package routes

import (
	"net/http"

	"github.com/akdevsaha-dev/applied-ratelimit-go/handlers"
)

func RegisterHomeRoute(mux *http.ServeMux) {
	mux.HandleFunc("/get-details", handlers.HomeHandler)
}
