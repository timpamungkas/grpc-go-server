package port

import (
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
)

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) float64
	FindExchangeRate(fromCur string, toCur string) float64
	CalculateTransactionSummary(tcur *dbank.TransactionSummary, trans dbank.Transaction) error
}
