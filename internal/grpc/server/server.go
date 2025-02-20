package server

import (
	"context"
	"fmt"
	pb "github.com/galkin09/proto-exchange/exchange"
	"google.golang.org/grpc"
	"gw-exchange/internal/storages"
	"gw-exchange/pkg/logs"
	"net"
)

type ExchangeServer struct {
	pb.UnimplementedExchangeServiceServer
	storage storages.Storage
	logger  *logs.Logger
}

func NewExchangeServer(storage storages.Storage, logger *logs.Logger) *ExchangeServer {
	return &ExchangeServer{
		storage: storage,
		logger:  logger.With("component", "grpc-server"),
	}
}

// GetExchangeRates возвращает все курсы валют
func (s *ExchangeServer) GetExchangeRates(ctx context.Context, _ *pb.Empty) (*pb.ExchangeRatesResponse, error) {
	s.logger.Info("Запрос на получение всех курсов валют")

	rates, err := s.storage.GetExchangeRates(ctx)
	if err != nil {
		s.logger.Error("Ошибка при получении курсов валют", err)
		return nil, fmt.Errorf("ошибка при получении курсов: %w", err)
	}

	response := &pb.ExchangeRatesResponse{
		Rates: rates,
	}
	return response, nil
}

// GetExchangeRateForCurrency возвращает курс для конкретной валютной пары
func (s *ExchangeServer) GetExchangeRateForCurrency(ctx context.Context, req *pb.CurrencyRequest) (*pb.ExchangeRateResponse, error) {
	s.logger.Info("Запрос на получение курса для валютной пары", "from", req.FromCurrency, "to", req.ToCurrency)

	rate, err := s.storage.GetExchangeRateForCurrency(ctx, req.FromCurrency, req.ToCurrency)
	if err != nil {
		s.logger.Error("Ошибка при получении курса для валютной пары", err)
		return nil, fmt.Errorf("ошибка при получении курса: %w", err)
	}

	response := &pb.ExchangeRateResponse{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         rate,
	}
	return response, nil
}

// Start запускает gRPC-сервер
func (s *ExchangeServer) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("не удалось запустить сервер: %w", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExchangeServiceServer(grpcServer, s)

	s.logger.Info("gRPC-сервер запущен", "addr", addr)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("ошибка при работе сервера: %w", err)
	}

	return nil
}
