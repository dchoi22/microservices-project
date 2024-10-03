package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dchoi22/microservices-project/logging-service/internal/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var serviceURL = os.Getenv("DATABASE_SERVICE_URL")
// This service uses the database-service API to fetch different log entries based on filters like status, time created, and endpoint id
func main() {

	if serviceURL == "" {
		log.Fatal("DATABASE_SERVICE_URL is not set")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/logs", func(r chi.Router) {
		r.Get("/", handlers.GetAllLogs)
		r.Get("/status/{status}", handlers.GetLogsByStatus)
		r.Get("/recent", handlers.GetRecentLogs)
		r.Get("/endpoint/{endpoint_id}", handlers.GetLogsByEndpointID)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}

}
