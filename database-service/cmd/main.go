package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dchoi22/microservices-project/database-service/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

// This service opens the PostgresDatabase and creates the endpoints and logs tables if they don't exist
func main() {

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Pings the db to verify that the connection is still active
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to the database")

	createTables(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := handlers.Handler{DB: db}

	// Endpoints CRUD endpoints
	r.Route("/endpoints", func(r chi.Router) {
		r.Get("/", h.GetEndpoints)
		r.Post("/", h.CreateEndpoint)
		r.Put("/{id}", h.UpdateEndpoint)
		r.Delete("/{id}", h.DeleteEndpoint)
	})

	// Logs CRUD endpoints, along with fetching logs by filters
	r.Route("/logs", func(r chi.Router) {
		r.Get("/", h.GetLogs)
		r.Get("/recent", h.GetRecentLogs)
		r.Get("/status/{status}", h.GetLogsByStatus)
		r.Get("/endpoint/{endpoint_id}", h.GetLog)
		r.Post("/", h.CreateLog)
	})

	// r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	fmt.Println("Health URL")
	// 	w.WriteHeader(http.StatusOK)
	// })

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// Function to create tables if they don't exist in the db
func createTables(db *sql.DB) {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS endpoints (
        id SERIAL PRIMARY KEY,
        url TEXT NOT NULL
    );
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS logs (
        id SERIAL PRIMARY KEY,
        endpoint_id INTEGER REFERENCES endpoints(id),
        status TEXT NOT NULL,
        timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
`)
	if err != nil {
		log.Fatal(err)
	}
}
