package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dchoi22/microservices-project/database-service/internal/models"
	"github.com/go-chi/chi/v5"
)
// Handler struct for accessing the db for operations
type Handler struct {
	DB *sql.DB
}
// Gets all of the endpoints in the endpoints table
func (h *Handler) GetEndpoints(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, url FROM endpoints")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var endpoints []models.Endpoint
	for rows.Next() {
		var e models.Endpoint
		if err := rows.Scan(&e.ID, &e.URL); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		endpoints = append(endpoints, e)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(endpoints); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Creates a new Endpoint with the request body data and inserts it into the endpoints table
func (h *Handler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var e models.Endpoint
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow("INSERT INTO endpoints (url) VALUES ($1) RETURNING id", e.URL).Scan(&e.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Updates the endpoint in the endpoints table specified by the id param 
func (h *Handler) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
		return
	}

	var e models.Endpoint
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec("UPDATE endpoints SET url = $1 WHERE id = $2", e.URL, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
// Deletes the endpoint in the endpoints table specified by the id param
func (h *Handler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, "Invalid endpoint ID", http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec("DELETE FROM endpoints WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
// Gets all of the logs in the logs table
func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, endpoint_id, status, timestamp FROM logs")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logEntries []models.LogEntry

	for rows.Next() {
		var l models.LogEntry
		if err := rows.Scan(&l.ID, &l.EndpointID, &l.Status, &l.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logEntries = append(logEntries, l)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Gets all of the logs in the logs table that were created less than a day ago
func (h *Handler) GetRecentLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`SELECT * FROM logs
		WHERE timestamp >= NOW() - INTERVAL '1 day'
		ORDER BY timestamp DESC;
`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var recentLogEntries []models.LogEntry

	for rows.Next() {
		var l models.LogEntry
		if err := rows.Scan(&l.ID, &l.EndpointID, &l.Status, &l.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		recentLogEntries = append(recentLogEntries, l)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(recentLogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Gets all of the logs in the logs table that were created with a specific status based on the status url param
func (h *Handler) GetLogsByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	rows, err := h.DB.Query(`SELECT * FROM logs WHERE status =$1`, status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	var statusLogEntries []models.LogEntry

	for rows.Next() {
		var l models.LogEntry
		if err := rows.Scan(&l.ID, &l.EndpointID, &l.Status, &l.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		statusLogEntries = append(statusLogEntries, l)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(statusLogEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Gets all of the logs in the logs table based on their endpoint id FKs
func (h *Handler) GetLog(w http.ResponseWriter, r *http.Request) {
	endpointID := chi.URLParam(r, "endpoint_id")
	rows, err := h.DB.Query(`SELECT * FROM logs WHERE endpoint_id =$1`, endpointID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logEntries []models.LogEntry

	for rows.Next() {
		var l models.LogEntry
		if err := rows.Scan(&l.ID, &l.EndpointID, &l.Status, &l.Timestamp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logEntries = append(logEntries, l)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(logEntries); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
// Creates a new Log with the request body data and inserts it into the logs table
func (h *Handler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var l models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.DB.QueryRow("INSERT INTO logs (endpoint_id, status) VALUES ($1, $2) RETURNING id, timestamp", l.EndpointID, l.Status).Scan(&l.ID, &l.Timestamp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(l); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
