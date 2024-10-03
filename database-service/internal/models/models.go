package models
// Endpoint data model
type Endpoint struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}
// Log entry data model 
type LogEntry struct {
	ID         int    `json:"id"`
	EndpointID int    `json:"endpoint_id"`
	Status     string `json:"status"`
	Timestamp  string `json:"timestamp"`
}
