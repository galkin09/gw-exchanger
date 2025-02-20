package storages

type ExchangeRate struct {
	CurrencyCode string  `json:"currency_code"` // Код валюты, например, USD
	Rate         float64 `json:"rate"`          // Курс валюты
}
