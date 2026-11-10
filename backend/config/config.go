package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置，来源为 .env 文件与系统环境变量
type Config struct {
	ServerHost string        // 服务监听地址
	ServerPort string        // 服务监听端口
	ServerMode string        // 运行模式：debug / release / test
	DBPath     string        // SQLite 数据库文件路径
	JWTSecret  string        // JWT 签名密钥
	JWTExpire  time.Duration // JWT 有效期
}

// Load 加载 .env（存在时）并读取配置，缺失项提供安全默认值
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("未找到 .env 文件，使用系统环境变量与默认值: %v", err)
	}

	return &Config{
		ServerHost: getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		ServerMode: getEnv("SERVER_MODE", "debug"),
		DBPath:     getEnv("DB_PATH", "blog.db"),
		JWTSecret:  getEnv("JWT_SECRET", "blog-app-secret-key-change-me"),
		JWTExpire:  getDuration("JWT_EXPIRE", 24*time.Hour),
	}
}

// getEnv 读取字符串环境变量，为空时返回默认值
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getDuration 读取时长环境变量（如 24h），非法时返回默认值
func getDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("环境变量 %s=%q 不是合法时长，使用默认值 %s", key, value, fallback)
		return fallback
	}
	return duration
}
