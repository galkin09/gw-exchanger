package storages

// Currency представляет модель для таблицы currencies
type Currency struct {
	Code      string  `json:"currency_code" db:"currency_code"` // Код валюты (например, USD, EUR)
	RateToUSD float32 `json:"rate_to_usd" db:"rate_to_usd"`     // Курс валюты к USD
}

// CurrencyRequest представляет запрос для получения курса валюты
type CurrencyRequest struct {
	From string `json:"from_currency"` // Исходная валюта
	To   string `json:"to_currency"`   // Целевая валюта
}

// CurrencyResponse представляет ответ с курсом валюты
type CurrencyResponse struct {
	From string  `json:"from_currency"` // Исходная валюта
	To   string  `json:"to_currency"`   // Целевая валюта
	Rate float32 `json:"rate"`          // Курс обмена
}
