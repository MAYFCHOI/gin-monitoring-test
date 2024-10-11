package cron

import (
	"github.com/robfig/cron/v3"
)

func SetupCron() (*cron.Cron, error) {
	c := cron.New()
	_, err := c.AddFunc("*/1 * * * *", fetchMetrics)
	if err != nil {
		return nil, err
	}

	return c, nil
}
