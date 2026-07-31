package config

import (
	"os"

	"github.com/ikhwan11/auth-be/internal/database"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppName string
	AppPort string

	EmployeeDB database.Config
	AuthDB     database.Config

	JWTSecret  string
	JWTExpired string
}

func Load() AppConfig {
	_ = godotenv.Load()

	return AppConfig{
		AppName: getEnv("APP_NAME", "auth-be"),
		AppPort: getEnv("APP_PORT", "8060"),

		EmployeeDB: database.Config{
			Host:     getEnv("EMPLOYEE_DB_HOST", "localhost"),
			Port:     getEnv("EMPLOYEE_DB_PORT", "5432"),
			User:     getEnv("EMPLOYEE_DB_USER", "postgres"),
			Password: getEnv("EMPLOYEE_DB_PASSWORD", ""),
			DBName:   getEnv("EMPLOYEE_DB_NAME", "employee"),
			SSLMode:  getEnv("EMPLOYEE_DB_SSLMODE", "disable"),
		},

		AuthDB: database.Config{
			Host:     getEnv("AUTH_DB_HOST", "localhost"),
			Port:     getEnv("AUTH_DB_PORT", "5432"),
			User:     getEnv("AUTH_DB_USER", "postgres"),
			Password: getEnv("AUTH_DB_PASSWORD", ""),
			DBName:   getEnv("AUTH_DB_NAME", "auth"),
			SSLMode:  getEnv("AUTH_DB_SSLMODE", "disable"),
		},

		JWTSecret:  getEnv("JWT_SECRET", ""),
		JWTExpired: getEnv("JWT_EXPIRED", "15m"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
