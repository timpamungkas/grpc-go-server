package bank

import "time"

const (
	Unknown string = "UNKNOWN"
	In      string = "IN"
	Out     string = "OUT"
)

type Account struct {
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	Transactions   []Transaction
}

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

type TransferRequest struct {
	FromAccountNumber string
	ToAccountNumber   string
	Currency          string
	Amount            float64
}
