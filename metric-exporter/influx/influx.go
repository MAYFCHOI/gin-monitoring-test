package influx

import (
	"context"
	"metric-exporter/config"
	"metric-exporter/model"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func InfluxInit() {
	config.DB = influxdb2.NewClient(config.Env.Database.Url, config.Env.Database.Token)
}

func GinMetricToInflux(metricsResponse model.MetricsResponse) {
	now := time.Now()
	writeAPI := config.DB.WriteAPIBlocking(config.Env.Database.Org, config.Env.Database.Bucket)

	for path, metric := range metricsResponse.Metrics {
		duration := metric.AvgDuration * float64(metric.Count)
		p := influxdb2.NewPointWithMeasurement(metricsResponse.ServiceName).
			AddTag("path", path).
			AddField("count", metric.Count).
			AddField("duration", duration)
		for k, v := range metric.StatusCodes {
			p.AddField(k, v)
		}
		p.SetTime(now)
		writeAPI.WritePoint(context.Background(), p)
	}
}
