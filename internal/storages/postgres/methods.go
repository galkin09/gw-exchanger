package postgres

import (
	"context"
	"fmt"
)

// GetExchangeRates Реализация интерфейса Storage
func (p *PSQL) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	const op = "GetExchangeRates psql"

	logger := p.logger.With("op", op)
	logger.Info("Запрос на получение курсов валют")

	rows, err := p.pool.Query(ctx, "SELECT currency_code, rate_to_usd FROM currencies")
	if err != nil {
		logger.Error("Ошибка выполнения запроса", err)
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	rates := make(map[string]float32)
	for rows.Next() {
		var currencyCode string
		var rateToUSD float32
		if err := rows.Scan(&currencyCode, &rateToUSD); err != nil {

			logger.Error("Ошибка сканирования строки", err) // Логируем ошибку
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rates[currencyCode] = rateToUSD
	}

	logger.Info("Запрос на получение курсов завершен успешно", "op", op)
	return rates, nil
}

func (p *PSQL) GetExchangeRateForCurrency(ctx context.Context, from, to string) (float32, error) {

	const op = "GetExchangeRateForCurrency"
	logger := p.logger.With("op", op)

	var fromRate, toRate float32

	// Получаем курс для обеих валют относительно USD
	err := p.pool.QueryRow(ctx, "SELECT rate_to_usd FROM currencies WHERE currency_code=$1", from).Scan(&fromRate)
	if err != nil {
		logger.Error("Ошибка выполнения запроса для валютной пары", err) // Логируем ошибку
		return 0, fmt.Errorf("%s: ошибка получения курса для %s: %w", op, from, err)
	}

	err = p.pool.QueryRow(ctx, "SELECT rate_to_usd FROM currencies WHERE currency_code=$1", to).Scan(&toRate)
	if err != nil {
		return 0, fmt.Errorf("%s: ошибка получения курса для %s: %w", op, to, err)
	}

	// Вычисляем курс между валютами
	rate := fromRate / toRate
	logger.Info("Запрос для валютной пары завершен успешно", "op", op)
	return rate, nil
}
