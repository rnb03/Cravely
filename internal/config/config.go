package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	ServerPort   int
	MongoURI     string
	DatabaseName string
}

// LoadConfig loads the configuration from environment variables
func LoadConfig() Config {
	config := Config{
		ServerPort:   8080,
		MongoURI:     "mongodb://localhost:27017",
		DatabaseName: "creavely",
	}

	if port := os.Getenv("SERVER_PORT"); port != "" {
		if val, err := strconv.Atoi(port); err == nil {
			config.ServerPort = val
		}
	}

	if uri := os.Getenv("MONGO_URI"); uri != "" {
		config.MongoURI = uri
	}

	if dbName := os.Getenv("DATABASE_NAME"); dbName != "" {
		config.DatabaseName = dbName
	}

	return config
}
