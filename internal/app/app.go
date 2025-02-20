package app

import (
	"context"
	"gw-exchange/internal/config"
	"gw-exchange/internal/storages/postgres"
	"gw-exchange/pkg/logs"
	"log"
	"time"
)

func Run() {
	const op = "GRPC Run"
	ctx := context.Background()

	cfg, err := config.LoadConfig("./config.env")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	storage := postgres.New()
	err = storage.Start(ctx, cfg.DBUrl, 5*time.Second, cfg.MigrationsPath)
	if err != nil {
		logs.Logger.Error("Ошибка подключения к БД", err)
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer storage.Stop()

	if err := StartGRPCServer(storage, cfg.GRPCPort); err != nil {
		logs.Logger.Error("Ошибка запуска gRPC сервера", err)
		log.Fatalf("Ошибка работы gRPC-сервера: %v", err)
	}
	logs.Logger.Info("Приложение успешно запущено", logs.Attr{Key: "operation", Value: op})
}
