package config

import (
	"os"
)

// 配置結構
type Config struct {
	Port string
}

// 加載配置
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9090" // 默認端口
	}

	return &Config{
		Port: port,
	}
}
