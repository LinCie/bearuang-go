package config

import (
	"errors"
	"os"
	"strconv"
)

const (
	defaultEnvironment    = "development"
	productionEnvironment = "production"
	defaultJWTSecret      = "development-secret"

	// AccessTokenCookieName is the name of the HTTP-only access-token cookie.
	AccessTokenCookieName = "bearuang_access_token"
	// RefreshTokenCookieName is the name of the HTTP-only refresh-token cookie.
	RefreshTokenCookieName = "bearuang_refresh_token"
)

type Config struct {
	Environment  string
	DatabaseURL  string
	Port         int
	JWTSecret    string
	CookieSecure bool
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

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = defaultEnvironment
	}

	return Config{
		Environment:  environment,
		DatabaseURL:  databaseURL,
		Port:         port,
		JWTSecret:    jwtSecret,
		CookieSecure: environment == productionEnvironment,
	}, nil
}
