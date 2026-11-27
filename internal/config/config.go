// Package config contains environment variables
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DBURL            string
	TemporalHostPort string
	RedisAddress     string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	DBURL = os.Getenv("DB_URL")
	TemporalHostPort = os.Getenv("TEMPORAL_HOST_PORT")
	RedisAddress = os.Getenv("REDIS_ADDRESS")
}
