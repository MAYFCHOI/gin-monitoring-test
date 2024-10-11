package cron

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"metric-exporter/config"
	"metric-exporter/influx"
	"metric-exporter/model"
	"net/http"
	"time"
)

func fetchMetrics() {
	for _, server := range config.Env.ServerInfo {
		metrics, err := fetchMetric(server)
		if err != nil {
			log.Printf("Error fetching metrics from server %s: %v\n", server.Url, err)
			continue
		}
		go influx.GinMetricToInflux(*metrics)
	}
}

func fetchMetric(server config.ServerInfo) (*model.MetricsResponse, error) {
	url := fmt.Sprintf("http://%s:%s%s", server.Url, server.Port, server.Path)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metrics from %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var metrics model.MetricsResponse
	if err := json.Unmarshal(body, &metrics); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return &metrics, nil
}
