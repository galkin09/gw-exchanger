package storages

import "context"

// Storage - интерфейс хранилища курсов валют
type Storage interface {
	GetExchangeRates(ctx context.Context) (map[string]float32, error)
	GetExchangeRateForCurrency(ctx context.Context, from, to string) (float32, error)
}
