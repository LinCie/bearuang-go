package config

import (
	"errors"
	"net"
	"net/url"
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
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		database := os.Getenv("POSTGRES_DB")
		host := os.Getenv("DATABASE_HOST")
		databasePort := os.Getenv("DATABASE_PORT")
		if user == "" || password == "" || database == "" || host == "" || databasePort == "" {
			return Config{}, errors.New("DATABASE_URL or all POSTGRES_* and DATABASE_HOST/DATABASE_PORT variables are required")
		}

		parsedPort, err := strconv.Atoi(databasePort)
		if err != nil || parsedPort <= 0 || parsedPort > 65535 {
			return Config{}, errors.New("DATABASE_PORT must be a valid port number")
		}

		dbURL := url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(user, password),
			Host:   net.JoinHostPort(host, databasePort),
			Path:   database,
		}
		query := dbURL.Query()
		query.Set("sslmode", "disable")
		dbURL.RawQuery = query.Encode()
		databaseURL = dbURL.String()
		if err := os.Setenv("DATABASE_URL", databaseURL); err != nil {
			return Config{}, err
		}
	}

	port := 8080
	if value := os.Getenv("PORT"); value != "" {
		parsedPort, err := strconv.Atoi(value)
		if err != nil || parsedPort <= 0 {
			return Config{}, errors.New("PORT must be a positive integer")
		}
		port = parsedPort
	}

	return Config{
		DatabaseURL: databaseURL,
		Port:        port,
	}, nil
}
