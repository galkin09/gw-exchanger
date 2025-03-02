DROP TABLE IF EXISTS currencies;

CREATE TABLE IF NOT EXISTS currencies  (
                                           currency_code VARCHAR(3) PRIMARY KEY,
                                           rate_to_usd FLOAT NOT NULL
);

INSERT INTO currencies (currency_code, rate_to_usd) VALUES
                                                        ('USD', 1.0),
                                                        ('EUR', 0.9),
                                                        ('RUB', 75.0);