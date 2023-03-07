package bank

import "time"

const (
	TransactionStatusUnknown string = "UNKNOWN"
	TransactionStatusIn      string = "IN"
	TransactionStatusOut     string = "OUT"
)

type Transaction struct {
	Amount          float64
	Timestamp       time.Time
	TransactionType string
	Notes           string
}

type ExchangeRate struct {
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimestamp time.Time
	ValidToTimestamp   time.Time
}

type TransactionSummary struct {
	SummaryOnDate time.Time
	SumIn         float64
	SumOut        float64
	SumTotal      float64
}
