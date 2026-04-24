package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL      string
	Port             string
	JWTAccessSecret  string
	JWTRefreshSecret string
	Env              string
}

func ExtractKey(key string, fallback string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}
	return value, nil
}

func LoadConfig() (*Config, error) {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
	}

	dbUrl, err := ExtractKey("DATABASE_URL", "")
	if err != nil {
		return nil, err
	}

	port, err := ExtractKey("PORT", "9000")
	if err != nil {
		return nil, err
	}

	jwtAccessSecret, err := ExtractKey("JWT_ACCESS_SECRET", "")
	if err != nil {
		return nil, err
	}

	jwtRefreshSecret, err := ExtractKey("JWT_REFRESH_SECRET", "")
	if err != nil {
		return nil, err
	}

	env, err := ExtractKey("ENV", "development")
	if err != nil {
		return nil, err
	}

	return &Config{
		DatabaseURL:      dbUrl,
		Port:             port,
		JWTAccessSecret:  jwtAccessSecret,
		JWTRefreshSecret: jwtRefreshSecret,
		Env:              env,
	}, nil
}