package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"gw-exchange/internal/config"
	"time"
)

type PSQL struct {
	pool    *pgxpool.Pool
	timeout time.Duration
	logger  *zap.Logger
}

// NewPSQL создает новый экземпляр PSQL
func NewPSQL(cfg config.PostgresConfig, logger *zap.Logger) (*PSQL, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	// Создаем пул соединений
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("ошибка при парсинге конфигурации пула: %w", err)
	}

	// Устанавливаем таймаут подключения
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ConnTimeout)*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании пула соединений: %w", err)
	}

	// Проверяем подключение
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ошибка при подключении к базе данных: %w", err)
	}

	logger.Info("Успешное подключение к PostgreSQL", zap.String("host", cfg.Host), zap.String("db", cfg.DBName))

	return &PSQL{
		pool:   pool,
		logger: logger,
	}, nil
}

func (p *PSQL) Start(ctx context.Context, url string, timeout time.Duration, migrationsPath string) error {
	const op = "PSQL START"
	p.timeout = timeout

	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	p.logger.Info("Подключаемся к базе данных", zap.String("url", url))

	pool, err := pgxpool.New(ctxTimeout, url)
	if err != nil {
		p.logger.Error("Ошибка при подключении к базе данных", zap.String("url", url), zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	p.pool = pool

	err = p.doMigrate(url, migrationsPath)
	if err != nil {
		p.logger.Error("Ошибка при выполнении миграции", zap.Error(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Info("Подключение и миграция успешно завершены")
	return nil
}

func (p *PSQL) doMigrate(dbURL, migrationsPath string) error {
	const op = "postgres.doMigrate"
	if migrationsPath == "" {
		return fmt.Errorf("%s: %s", op, "путь к миграциям пуст")
	}

	m, err := migrate.New("file://"+migrationsPath, dbURL)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.logger.Info("Миграции успешно применены", zap.String("path", migrationsPath))
	return nil
}

func (p *PSQL) Stop() {
	if p.pool != nil {
		p.logger.Info("Закрытие пула соединений с базой данных")
		p.pool.Close()
		p.logger.Info("Пул соединений успешно закрыт")
	} else {
		p.logger.Warn("Пул соединений уже закрыт или не был инициализирован")
	}
}
