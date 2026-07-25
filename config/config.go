package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	Name                   string
	SSLMode                string
	MaxConnections         int32
	MinConnections         int32
	MaxConnLifetime        time.Duration
	MaxConnIdleTime        time.Duration
	HealthCheckPeriod      time.Duration
}

func LoadConfig() (*Config, error) {
	maxConnections, err := getEnvInt32("DB_MAX_CONNECTIONS", 20)
	if err != nil {
		return nil, err
	}

	minConnections, err := getEnvInt32("DB_MIN_CONNECTIONS", 2)
	if err != nil {
		return nil, err
	}

	maxConnLifetimeMinutes, err := getEnvInt(
		"DB_MAX_CONN_LIFETIME_MINUTES",
		30,
	)
	if err != nil {
		return nil, err
	}

	maxConnIdleTimeMinutes, err := getEnvInt(
		"DB_MAX_CONN_IDLE_TIME_MINUTES",
		10,
	)
	if err != nil {
		return nil, err
	}

	healthCheckSeconds, err := getEnvInt(
		"DB_HEALTH_CHECK_SECONDS",
		30,
	)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "soto-ordering"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:           getEnv("DB_HOST", "localhost"),
			Port:           getEnv("DB_PORT", "5432"),
			User:           getEnv("DB_USER", "postgres"),
			Password:       getEnv("DB_PASSWORD", ""),
			Name:           getEnv("DB_NAME", "soto_ordering"),
			SSLMode:        getEnv("DB_SSLMODE", "disable"),
			MaxConnections: maxConnections,
			MinConnections: minConnections,

			MaxConnLifetime: time.Duration(maxConnLifetimeMinutes) *
				time.Minute,

			MaxConnIdleTime: time.Duration(maxConnIdleTimeMinutes) *
				time.Minute,

			HealthCheckPeriod: time.Duration(healthCheckSeconds) *
				time.Second,
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST tidak boleh kosong")
	}

	if c.Database.Port == "" {
		return fmt.Errorf("DB_PORT tidak boleh kosong")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER tidak boleh kosong")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME tidak boleh kosong")
	}

	if c.Database.MaxConnections < 1 {
		return fmt.Errorf("DB_MAX_CONNECTIONS minimal 1")
	}

	if c.Database.MinConnections < 0 {
		return fmt.Errorf("DB_MIN_CONNECTIONS tidak boleh negatif")
	}

	if c.Database.MinConnections > c.Database.MaxConnections {
		return fmt.Errorf(
			"DB_MIN_CONNECTIONS tidak boleh lebih besar dari DB_MAX_CONNECTIONS",
		)
	}

	return nil
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Port,
		d.Name,
		d.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(key string, fallback int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf(
			"environment variable %s harus berupa angka: %w",
			key,
			err,
		)
	}

	return parsedValue, nil
}

func getEnvInt32(key string, fallback int32) (int32, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf(
			"environment variable %s harus berupa angka: %w",
			key,
			err,
		)
	}

	return int32(parsedValue), nil
}