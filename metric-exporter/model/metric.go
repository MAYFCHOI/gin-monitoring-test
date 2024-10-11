package model

type MetricsResponse struct {
	Metrics     map[string]MetricDetail `json:"metrics"`
	ServiceName string                  `json:"service_name"`
}

type MetricDetail struct {
	AvgDuration float64        `json:"avg_duration"`
	Count       int            `json:"count"`
	StatusCodes map[string]int `json:"status_codes"`
}
