package models

// HealthResp struct
type HealthResp struct {
	ServiceURL string `json:"service_url"`
	Status     string `json:"status"`
	Error      string `json:"error,omitempty"`
}
