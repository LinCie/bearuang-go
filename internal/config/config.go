package config

import (
	"errors"
	"os"
	"strconv"
)

const defaultJWTSecret = "development-secret"

type Config struct {
	DatabaseURL string
	Port        int
	JWTSecret   string
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

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
		JWTSecret:   jwtSecret,
	}, nil
}
