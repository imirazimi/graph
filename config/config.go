package config

import (
    "fmt"
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    AppPort string
    AppEnv  string

    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("warning: .env file not found")
    }

    return &Config{
        AppPort: getEnv("APP_PORT", "8080"),
        AppEnv:  getEnv("APP_ENV", "development"),

        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "task_manager"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
    }
}

func (c *Config) DatabaseURL() string {
    return fmt.Sprintf(
        "postgres://%s:%s@%s:%s/%s?sslmode=%s",
        c.DBUser,
        c.DBPassword,
        c.DBHost,
        c.DBPort,
        c.DBName,
        c.DBSSLMode,
    )
}

func getEnv(key, fallback string) string {
    value := os.Getenv(key)
    if value == "" {
        return fallback
    }

    return value
}