package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dchoi22/microservices-project/health-check-service/internal/models"
)

var serviceURL = os.Getenv("DATABASE_SERVICE_URL")

// This service uses the Database-service API to test the endpoints found in the enpoints table and returns their status every minute
func main() {

	if serviceURL == "" {
		log.Fatal("DATABASE_SERVICE_URL is not set")
	}

	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		performHealthCheck()
	}
	// r := chi.NewRouter()
	// fmt.Print(r)
}

func performHealthCheck() {
	endPoints, err := GetEndpoints()
	if err != nil {
		log.Printf("Error getting endpoints: %v", err)
		return
	}
	for _, endpoint := range endPoints {
		status := checkEndpoint(endpoint.URL)
		logResult(endpoint.ID, status)
	}
}

func checkEndpoint(url string) string {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "DOWN"
	}
	return "OK"
}

func GetEndpoints() ([]models.Endpoint, error) {
	resp, err := http.Get(fmt.Sprintf("%s/endpoints", serviceURL))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get endpoints: status code %d", resp.StatusCode)
	}

	var endpoints []models.Endpoint
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, fmt.Errorf("failed to get endpoints: status code %d", resp.StatusCode)
	}
	return endpoints, nil
}

func logResult(id int, status string) {
	logEntry := models.LogEntry{
		EndpointID: id,
		Status:     status,
	}

	data, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf("Error marshaling log entry: %v", err)
		return
	}

	resp, err := http.Post(fmt.Sprintf("%s/logs", serviceURL), "application/json", bytes.NewReader(data))
	if err != nil {
		log.Printf("Error logging result for endpoint %d: %v", id, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to log result for endpoint %d: status code %d", id, resp.StatusCode)
	}
}
