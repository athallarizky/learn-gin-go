package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPass             string
	DBName             string
	DBSSL              string
	GithubClientID     string
	GithubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
	AuthRedirectURL    string
}

var App *Config

func init() {
	LoadConfig()
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("Note: .env file not found, using system env/defaults")
	}

	App = &Config{
		AppPort:            getEnv("PORT", "8088"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5433"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPass:             getEnv("DB_PASSWORD", "postgres"),
		DBName:             getEnv("DB_NAME", "learn_gin_db"),
		DBSSL:              getEnv("DB_SSLMODE", "disable"),
		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		AuthRedirectURL:    getEnv("AUTH_REDIRECT_URL", "http://localhost:8088/api/v1/auth"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
