package config

import (
	"os"
	"time"
)

// Config 系统配置
type Config struct {
	APIKey       string
	DatabaseURL  string
	Host         string
	Port         string
	TicketPrefix string

	// JWT 鉴权（管理后台用）
	JWTSecret   string
	JWTExpiry   time.Duration
	AdminUser   string
	AdminPass   string
}

// Load 从环境变量加载配置，提供默认值
func Load() *Config {
	return &Config{
		APIKey:       getEnv("API_KEY", "sk-change-me-in-production"),
		DatabaseURL:  getEnv("DATABASE_URL", "tickets.db"),
		Host:         getEnv("HOST", "0.0.0.0"),
		Port:         getEnv("PORT", "8000"),
		TicketPrefix: getEnv("TICKET_PREFIX", "T"),

		// JWT 配置
		JWTSecret: getEnv("JWT_SECRET", "ticket-system-dev-secret-change-me-in-production"),
		JWTExpiry: parseDuration(getEnv("JWT_EXPIRY", "8h")),
		AdminUser: getEnv("ADMIN_USERNAME", "admin"),
		AdminPass: getEnv("ADMIN_PASSWORD", "admin123"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 8 * time.Hour
	}
	return d
}
