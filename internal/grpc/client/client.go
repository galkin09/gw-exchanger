package client

import (
	"context"
	"fmt"
	pb "github.com/galkin09/proto-exchange/exchange"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gw-exchange/pkg/logs"
)

type ExchangeClient struct {
	client pb.ExchangeServiceClient
	logger *logs.Logger
}

func NewExchangeClient(addr string, logger *logs.Logger) (*ExchangeClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к серверу: %w", err)
	}

	client := pb.NewExchangeServiceClient(conn)
	return &ExchangeClient{
		client: client,
		logger: logger.With("component", "grpc-client"),
	}, nil
}

// GetExchangeRates получает все курсы валют
func (c *ExchangeClient) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	c.logger.Info("Запрос на получение всех курсов валют")

	response, err := c.client.GetExchangeRates(ctx, &pb.Empty{})
	if err != nil {
		c.logger.Error("Ошибка при получении курсов валют", err)
		return nil, fmt.Errorf("ошибка при получении курсов: %w", err)
	}

	return response.Rates, nil
}

// GetExchangeRateForCurrency получает курс для конкретной валютной пары
func (c *ExchangeClient) GetExchangeRateForCurrency(ctx context.Context, from, to string) (float32, error) {
	c.logger.Info("Запрос на получение курса для валютной пары", "from", from, "to", to)

	response, err := c.client.GetExchangeRateForCurrency(ctx, &pb.CurrencyRequest{
		FromCurrency: from,
		ToCurrency:   to,
	})
	if err != nil {
		c.logger.Error("Ошибка при получении курса для валютной пары", err)
		return 0, fmt.Errorf("ошибка при получении курса: %w", err)
	}

	return response.Rate, nil
}
