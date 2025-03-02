package postgres

import (
	"context"
	"fmt"
	"go.uber.org/zap"
)

// GetExchangeRates Реализация интерфейса Storage
func (p *PSQL) GetExchangeRates(ctx context.Context) (map[string]float32, error) {
	const op = "GetExchangeRates psql"

	p.logger.Info("Запрос на получение курсов валют")

	rows, err := p.pool.Query(ctx, "SELECT currency_code, rate_to_usd FROM currencies")
	if err != nil {
		p.logger.Error("Ошибка выполнения запроса", zap.Error(err))
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	rates := make(map[string]float32)
	for rows.Next() {
		var currencyCode string
		var rateToUSD float32
		if err := rows.Scan(&currencyCode, &rateToUSD); err != nil {

			p.logger.Error("Ошибка сканирования строки", zap.Error(err)) // Логируем ошибку
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		rates[currencyCode] = rateToUSD
	}

	p.logger.Info("Запрос на получение курсов завершен успешно", zap.String("op", op))
	return rates, nil
}

// GetExchangeRateForCurrency Реализация интерфейса Storage
func (p *PSQL) GetExchangeRateForCurrency(ctx context.Context, from, to string) (float32, error) {

	const op = "GetExchangeRateForCurrency"

	var fromRate, toRate float32

	// Используем один запрос для получения курсов обеих валют
	query := `
        SELECT 
            MAX(CASE WHEN currency_code = $1 THEN rate_to_usd END) AS from_rate,
            MAX(CASE WHEN currency_code = $2 THEN rate_to_usd END) AS to_rate
        FROM currencies
        WHERE currency_code IN ($1, $2)
    `

	err := p.pool.QueryRow(ctx, query, from, to).Scan(&fromRate, &toRate)
	if err != nil {
		p.logger.Error("Ошибка выполнения запроса для валютной пары", zap.Error(err))
		return 0, fmt.Errorf("%s: ошибка получения курса для пары %s/%s: %w", op, from, to, err)
	}

	// Проверяем, что курсы были найдены
	if fromRate == 0 || toRate == 0 {
		p.logger.Error("Курс для одной из валют не найден", zap.Error(err))
		return 0, fmt.Errorf("%s: курс для одной из валют (%s/%s) не найден", op, from, to)
	}

	// Вычисляем курс между валютами
	rate := toRate / fromRate
	p.logger.Info("Запрос для валютной пары завершен успешно", zap.String("op", op))
	return rate, nil
}
