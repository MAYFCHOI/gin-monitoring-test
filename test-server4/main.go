package main

import (
	"log"
	"net/http"

	"github.com/MAYFCHOI/gin-monitoring/metrics"
	"github.com/MAYFCHOI/gin-monitoring/tracing"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(metrics.MetricsMiddleware(
		metrics.MetricInit{
			ServiceName: "test-server4"}))

	r.Use(tracing.TracingMiddleware(
		tracing.TraceInit{
			ServiceName: "test-server4",
			Logpath:     "../logs/test-server4.log"}))

	client := &http.Client{
		Transport: tracing.NewTracingTransport(http.DefaultTransport),
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/metrics", metrics.MetricsHandler)
	r.GET("/process", processRequest(client))

	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func processRequest(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"})
	}
}
