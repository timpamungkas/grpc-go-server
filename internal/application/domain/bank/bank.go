package domainbank

import "time"

const (
	Unknown int32 = 0
	In      int32 = 1
	Out     int32 = -1
)

type Account struct {
	AccountNumber  string
	AccountName    string
	CurrentBalance float64
	Transactions   []Transaction
}

type Transaction struct {
	Amount          float64
	Timestamp       time.Time
	TransactionType int32
	Notes           string
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
