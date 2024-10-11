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
			ServiceName: "test-server1"}))

	r.Use(tracing.TracingMiddleware(
		tracing.TraceInit{
			ServiceName: "test-server1",
			Logpath:     "../logs/test-server.log"}))

	client := &http.Client{
		Transport: tracing.NewTracingTransport(http.DefaultTransport),
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	r.GET("/metrics", metrics.MetricsHandler)
	r.GET("/process", processRequest(client))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func processRequest(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := http.NewRequestWithContext(c.Request.Context(), "GET", "http://localhost:8081/process", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Backend service call failed"})
			return
		}
		defer resp.Body.Close()

		req2, err := http.NewRequestWithContext(c.Request.Context(), "GET", "http://localhost:8082/process", nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
			return
		}
		resp2, err := client.Do(req2)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Backend service call failed"})
			return
		}
		defer resp2.Body.Close()
		c.JSON(http.StatusOK, gin.H{"message": "Request processed successfully"})
	}
}
