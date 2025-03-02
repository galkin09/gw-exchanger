package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"os"
	"strconv"
)

// Config представляет конфигурацию приложения
type Config struct {
	Postgres PostgresConfig
	GRPC     GRPCConfig
	Logger   *zap.Logger
}

// PostgresConfig содержит параметры для подключения к базе данных
type PostgresConfig struct {
	Host           string
	Port           string
	User           string
	Password       string
	DBName         string
	ConnTimeout    int
	MigrationsPath string
}

// GRPCConfig содержит параметры для настройки gRPC сервера
type GRPCConfig struct {
	Host    string
	Port    string
	Timeout int
}

// ConnectionURL генерирует строку подключения к PostgreSQL
func (c *PostgresConfig) ConnectionURL() (string, error) {
	if c.Host == "" || c.Port == "" || c.DBName == "" || c.User == "" || c.Password == "" {
		return "", fmt.Errorf("некоторые параметры подключения отсутствуют")
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DBName), nil
}

// ConnectionURL генерирует строку подключения для gRPC
func (g GRPCConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%s", g.Host, g.Port)
}

// LoadConfig загружает конфигурацию из файла .env.
func LoadConfig(envPath string) (*Config, error) {
	const op = "config.LoadConfig"

	if err := godotenv.Load(envPath); err != nil {
		return nil, fmt.Errorf("%s: ошибка загрузки .env файла: %w", op, err)
	}

	connTimeout, err := strconv.Atoi(os.Getenv("DB_CONN_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: неверное значение для DB_CONN_TIMEOUT: %w", op, err)
	}

	postgresConfig := PostgresConfig{
		User:           os.Getenv("DB_USER"),
		Password:       os.Getenv("DB_PASSWORD"),
		Host:           os.Getenv("DB_HOST"),
		Port:           os.Getenv("DB_PORT"),
		DBName:         os.Getenv("DB_NAME"),
		ConnTimeout:    connTimeout,
		MigrationsPath: os.Getenv("MIGRATIONS_PATH"),
	}

	connGRPCTimeout, err := strconv.Atoi(os.Getenv("GRPC_TIMEOUT"))
	grpcConfig := GRPCConfig{
		Host:    os.Getenv("GRPC_HOST"),
		Port:    os.Getenv("GRPC_PORT"),
		Timeout: connGRPCTimeout,
	}

	return &Config{
		Postgres: postgresConfig,
		GRPC:     grpcConfig,
	}, nil
}
