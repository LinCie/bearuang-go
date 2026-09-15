package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL string
	Port        int
}

func NewConfig() (Config, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	port := 8080
	if value := os.Getenv("PORT"); value != "" {
		parsedPort, err := strconv.Atoi(value)
		if err != nil {
			return Config{}, err
		}
		port = parsedPort
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}
