package config

import (
	"log"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/spf13/viper"
)

type Config struct {
	ServerInfo []ServerInfo
	Database   Database
}

type ServerInfo struct {
	Url  string
	Port string
	Path string
	Name string
}

type Database struct {
	Url    string
	Token  string
	Org    string
	Bucket string
}

var Env Config
var DB influxdb2.Client

func GetEnvironmentVariable() error {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("config.json")
	var config Config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
	Env = config

	return nil
}
