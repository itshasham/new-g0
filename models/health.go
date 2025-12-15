package models

// Health represents the JSON response emitted by /healthz.
type Health struct {
	Status string `json:"status"`
}
