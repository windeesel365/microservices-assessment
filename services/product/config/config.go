package config

import (
	"log"
	"os"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func LoadConfig() Config {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not set in the environment")
	}

	dbSource := os.Getenv("DATABASE_URL")
	if dbSource == "" {
		log.Fatalf("DATABASE_URL is not set in the environment")
	}

	return Config{
		Port:        port,
		DatabaseURL: dbSource,
	}
}
