package main

import (
	"log"
	"net/http"

	"github.com/akdevsaha-dev/applied-ratelimit-go/routes"
)

func main() {

	mux := http.NewServeMux()

	routes.RegisterHomeRoute(mux)

	log.Println("Server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed!", err)
	}
}
