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
	JWT      JWTConfig
	Midtrans	MidtransConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type JWTConfig struct {
	AccessSecret     string
	Issuer           string
	AccessExpiration time.Duration
}

type DatabaseConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	SSLMode           string
	MaxConnections    int32
	MinConnections    int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

func LoadConfig() (*Config, error) {
	accessExpirationMinutes, err := getEnvInt(
		"JWT_ACCESS_EXPIRE_MINUTES",
		15,
	)
	if err != nil {
		return nil, err
	}

	maxConnections, err := getEnvInt32(
		"DB_MAX_CONNECTIONS",
		20,
	)
	if err != nil {
		return nil, err
	}

	minConnections, err := getEnvInt32(
		"DB_MIN_CONNECTIONS",
		2,
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

		JWT: JWTConfig{
			AccessSecret: getEnv(
				"JWT_ACCESS_SECRET",
				"",
			),
			Issuer: getEnv(
				"JWT_ISSUER",
				"soto-ordering-api",
			),
			AccessExpiration: time.Duration(
				accessExpirationMinutes,
			) * time.Minute,
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

			MaxConnLifetime:   30 * time.Minute,
			MaxConnIdleTime:   10 * time.Minute,
			HealthCheckPeriod: 30 * time.Second,
		},

		Midtrans: MidtransConfig{
			ServerKey: getEnv(
				"MIDTRANS_SERVER_KEY",
				"",
			),
			ClientKey: getEnv(
				"MIDTRANS_CLIENT_KEY",
				"",
			),
			IsProduction: getEnv(
				"MIDTRANS_IS_PRODUCTION",
				"false",
			) == "true",
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

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER tidak boleh kosong")
	}

	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME tidak boleh kosong")
	}

	if len(c.JWT.AccessSecret) < 32 {
		return fmt.Errorf(
			"JWT_ACCESS_SECRET minimal 32 karakter",
		)
	}

	if c.Midtrans.ServerKey == "" {
		return fmt.Errorf(
			"MIDTRANS_SERVER_KEY tidak boleh kosong",
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
			"%s harus berupa angka: %w",
			key,
			err,
		)
	}

	return parsedValue, nil
}

func getEnvInt32(
	key string,
	fallback int32,
) (int32, error) {
	value := os.Getenv(key)
	if value == "" {
		return fallback, nil
	}

	parsedValue, err := strconv.ParseInt(
		value,
		10,
		32,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"%s harus berupa angka: %w",
			key,
			err,
		)
	}

	return int32(parsedValue), nil
}