package app

import (
	"context"
	"gw-exchange/internal/config"
	"gw-exchange/internal/grpc/server"
	"gw-exchange/internal/storages/postgres"
	"gw-exchange/pkg/logs"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	const op = "GRPC Run"

	// Инициализация логгера
	logger := logs.New(os.Stdout)

	// Загрузка конфигурации из config.env
	if err := config.LoadConfig("config.env", logger); err != nil {
		logger.Error("Ошибка загрузки конфигурации", err)
		os.Exit(1)
	}

	// Создание конфигурации
	cfg := config.New(logger)

	// Инициализация хранилища
	storage, err := postgres.NewPSQL(cfg.Postgres, logger)
	if err != nil {
		panic(err)
	}

	dbURL, err := cfg.Postgres.ConnectionURL()
	if err != nil {
		logger.Error("Ошибка при формировании строки подключения к базе данных", err)
		os.Exit(1)
	}

	if err := storage.Start(context.Background(), dbURL, time.Duration(cfg.Postgres.ConnTimeout)*time.Second, cfg.MigrationsPath); err != nil {
		logger.Error("Ошибка при инициализации хранилища", err)
		os.Exit(1)
	}
	defer storage.Stop()

	// Запуск gRPC-сервера
	grpcServer := server.NewExchangeServer(storage, logger)
	go func() {
		if err := grpcServer.Start(cfg.GRPC.ConnectionURL()); err != nil {
			logger.Error("Ошибка при запуске gRPC-сервера", err)
			os.Exit(1)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Остановка сервера...")
}
