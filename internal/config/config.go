package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env       string
	Port      int
	AppSecret string

	DB     DBConfig
	Redis  RedisConfig
	Gate   GateConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type GateConfig struct {
	APIKey    string
	APISecret string
	BaseURL   string
	WSURL     string
}

var C *Config

func Init() error {
	// Load .env file
	_ = godotenv.Load()

	C = &Config{
		Env:       getEnv("APP_ENV", "development"),
		Port:      getEnvInt("APP_PORT", 8080),
		AppSecret: getEnv("APP_SECRET", "change-me-in-production"),

		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "gate_trading"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
		},

		Gate: GateConfig{
			APIKey:    getEnv("GATE_API_KEY", ""),
			APISecret: getEnv("GATE_API_SECRET", ""),
			BaseURL:   getEnv("GATE_API_BASE_URL", "https://api.gateio.ws/api/v4"),
			WSURL:     getEnv("GATE_WS_URL", "wss://api.gateio.ws/ws/v4/"),
		},
	}

	return nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
