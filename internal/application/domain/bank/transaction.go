package domainbank

import "time"

const (
	Unknown int32 = 0
	In      int32 = 1
	Out     int32 = -1
)

type Transaction struct {
	Amount          float64
	Timestamp       time.Time
	TransactionType int32
}

type TransactionSummary struct {
	SummaryOnDate time.Time
	SumIn         float64
	SumOut        float64
	SumTotal      float64
}
