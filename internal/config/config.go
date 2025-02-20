package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"

	"gw-exchange/pkg/logs"
)

// Config представляет конфигурацию приложения
type Config struct {
	Postgres       PostgresConfig
	GRPC           GRPCConfig
	Logger         logs.Logger
	MigrationsPath string
}

// PostgresConfig содержит параметры для подключения к базе данных
type PostgresConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	ConnTimeout int
}

// GRPCConfig содержит параметры для настройки gRPC сервера
type GRPCConfig struct {
	Host    string
	Port    string
	Timeout int
}

// New создает новый объект конфигурации
func New(logger *logs.Logger) *Config {
	cfg := &Config{
		Logger: *logger,
	}
	cfg.loadPostgresConfig()
	cfg.loadGRPCConfig()
	cfg.loadMigrationsPath()
	return cfg
}

// LoadConfig загружает конфигурацию из файла .env
func LoadConfig(configPath string, logger *logs.Logger) error {
	err := godotenv.Load(configPath)
	if err != nil {
		return fmt.Errorf("не удалось загрузить файл конфигурации: %w", err)
	}

	logger.Info("Конфигурация загружена из файла", "file", configPath)
	return nil
}

// loadPostgresConfig загружает конфигурацию для PostgreSQL из переменных окружения
func (c *Config) loadPostgresConfig() {
	c.Postgres.Host = getEnv("POSTGRES_HOST", "localhost")
	c.Postgres.Port = getEnv("POSTGRES_PORT", "5432")
	c.Postgres.User = getEnv("POSTGRES_USER", "postgres")
	c.Postgres.Password = getEnv("POSTGRES_PASSWORD", "postgres")
	c.Postgres.DBName = getEnv("POSTGRES_DB", "exchange_db")
	c.Postgres.ConnTimeout = getEnvAsInt("POSTGRES_CONN_TIMEOUT", 5)

	c.Logger.Info("Конфигурация PostgreSQL загружена", "host", c.Postgres.Host)
}

// ConnectionURL генерирует строку подключения к PostgreSQL
func (c *PostgresConfig) ConnectionURL() (string, error) {
	if c.Host == "" || c.Port == "" || c.User == "" || c.DBName == "" || c.Password == "" {
		return "", fmt.Errorf("некоторые параметры подключения отсутствуют")
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", c.User, c.Password, c.Host, c.Port, c.DBName), nil
}

// loadGRPCConfig загружает конфигурацию для gRPC сервера из переменных окружения
func (c *Config) loadGRPCConfig() {
	c.GRPC.Host = getEnv("GRPC_HOST", "localhost")
	c.GRPC.Port = getEnv("GRPC_PORT", "9090")
	c.GRPC.Timeout = getEnvAsInt("GRPC_TIMEOUT", 5)

	c.Logger.Info("Конфигурация gRPC загружена", "host", c.GRPC.Host)
}

// ConnectionURL генерирует строку подключения для gRPC
func (g GRPCConfig) ConnectionURL() string {
	return fmt.Sprintf("%s:%s", g.Host, g.Port)
}

// loadMigrationsPath загружает путь к миграциям из .env
func (c *Config) loadMigrationsPath() {
	c.MigrationsPath = getEnv("MIGRATIONS_PATH", "./migrations")
	c.Logger.Info("Путь к миграциям загружен", "path", c.MigrationsPath)
}

// getEnv возвращает значение переменной окружения или дефолтное значение
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt возвращает значение переменной окружения как целое число или дефолтное значение
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
