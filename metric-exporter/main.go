package main

import (
	"log"
	"metric-exporter/config"
	"metric-exporter/cron"
	"metric-exporter/gin"
	"metric-exporter/influx"
)

func main() {
	err := config.GetEnvironmentVariable()
	if err != nil {
		log.Panic(err)
	}

	influx.InfluxInit()

	c, err := cron.SetupCron()
	if err != nil {
		log.Println(err)
	}
	c.Start()

	r := gin.SetupRouter()
	err = r.Run(":9000")
	if err != nil {
		log.Panic(err)

	}
}
