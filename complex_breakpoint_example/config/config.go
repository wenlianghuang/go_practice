package config

import (
	"fmt"
	"os"
)

// 配置結構
type Config struct {
	Port         string
	DatabaseURL  string
	DatabaseHost string
	DatabasePort string
	DatabaseUser string
	DatabasePass string
	DatabaseName string
}

// 加載配置
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090" // 默認端口
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// 如果沒有提供 DATABASE_URL，從其他環境變量構建
		dbHost := getEnvOrDefault("DB_HOST", "localhost")
		dbPort := getEnvOrDefault("DB_PORT", "5432")
		dbUser := getEnvOrDefault("DB_USER", "postgres")
		dbPass := getEnvOrDefault("DB_PASSWORD", "postgres")
		dbName := getEnvOrDefault("DB_NAME", "breakpoint_db")

		databaseURL = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPass, dbName)
	}

	return &Config{
		Port:        port,
		DatabaseURL: databaseURL,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
