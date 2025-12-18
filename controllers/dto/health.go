package dto

// HealthResponse is the DTO returned by the GET /healthz endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}
