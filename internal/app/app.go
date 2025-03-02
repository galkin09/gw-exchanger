package app

import (
	"context"
	"go.uber.org/zap"
	"gw-exchange/internal/config"
	"gw-exchange/internal/grpc/server"
	"gw-exchange/internal/storages/postgres"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {
	const op = "GRPC Run"

	// Инициализация логгера
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.env")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}
	logger.Info("Конфигурация загружена", zap.Any("cfg", cfg))

	// Инициализация хранилища
	storage, err := postgres.NewPSQL(cfg.Postgres, logger)
	if err != nil {
		panic(err)
	}

	dbURL, err := cfg.Postgres.ConnectionURL()
	if err != nil {
		logger.Error("Ошибка при формировании строки подключения к базе данных", zap.Error(err))
		os.Exit(1)
	}

	if err := storage.Start(context.Background(), dbURL, time.Duration(cfg.Postgres.ConnTimeout)*time.Second, cfg.Postgres.MigrationsPath); err != nil {
		logger.Error("Ошибка при инициализации хранилища", zap.Error(err))
		os.Exit(1)
	}
	defer storage.Stop()

	// Запуск gRPC-сервера
	grpcServer := server.NewExchangeServer(storage, logger)
	go func() {
		if err := grpcServer.Start(cfg.GRPC.ConnectionURL()); err != nil {
			logger.Error("Ошибка при запуске gRPC-сервера", zap.Error(err))
			os.Exit(1)
		}
	}()

	// Ожидание сигналов для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Остановка сервера...")
}
