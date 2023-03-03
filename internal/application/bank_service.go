package application

import (
	"errors"
	"time"

	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
)

var accounts map[string]float64

type BankService struct {
}

func init() {
	accounts = map[string]float64{
		"111": 5001,
		"222": 5002,
		"333": 5003,
	}
}

func (a *BankService) FindCurrentBalance(acct string) float64 {
	return accounts[acct]
}

func (a *BankService) FindExchangeRate(fromCur string, toCur string) float64 {
	now := time.Now()
	bal := 1000 + now.Minute() + now.Second()

	return float64(bal)
}

func (a *BankService) CalculateTransactionSummary(tcur *dbank.TransactionSummary, tnew dbank.Transaction) error {
	switch tnew.TransactionType {
	case dbank.In:
		tcur.SumIn += tnew.Amount
	case dbank.Out:
		tcur.SumOut += tnew.Amount
	default:
		return errors.New("unknown transaction type")
	}

	tcur.SumTotal = tcur.SumIn - tcur.SumOut

	return nil
}

func (a *BankService) Transfer(fromAcct string, toAcct string, amount float64) (bool, error) {
	return true, nil
}
