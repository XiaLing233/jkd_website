package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBRUser    string
	DBRPass    string
	DBMeta     string
	ServerPort string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "mysql"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBRUser:    getEnv("DB_R_USER", "cirno"),
		DBRPass:    getEnv("DB_R_PASSWORD", ""),
		DBMeta:     getEnv("DB_META", "course_scheduler_meta"),
		ServerPort: getEnv("SERVER_PORT", "8000"),
	}

	if cfg.DBRPass == "" {
		return nil, fmt.Errorf("DB_R_PASSWORD 未设置")
	}
	return cfg, nil
}

func (c *Config) MetaDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		c.DBRUser, c.DBRPass, c.DBHost, c.DBPort, c.DBMeta,
	)
}

func (c *Config) CalendarDSN(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		c.DBRUser, c.DBRPass, c.DBHost, c.DBPort, dbName,
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
