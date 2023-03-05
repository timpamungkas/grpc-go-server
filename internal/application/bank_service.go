package application

import (
	"errors"
	"log"
	"time"

	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
	"github.com/timpamungkas/grpc-go-server/internal/port"
)

var accounts map[string]float64

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}

func init() {
	accounts = map[string]float64{
		"111": 5001,
		"222": 5002,
		"333": 5003,
	}
}

func (b *BankService) FindCurrentBalance(acct string) float64 {
	bankAccount, err := b.db.GetBankAccountByAccountNumber(acct, false)

	if err != nil {
		log.Printf("Error on FindCurrentBalance : %v\n", err)
	}

	return bankAccount.CurrentBalance
}

func (b *BankService) FindExchangeRate(fromCur string, toCur string) float64 {
	now := time.Now()
	bal := 1000 + now.Minute() + now.Second()

	return float64(bal)
}

func (b *BankService) CalculateTransactionSummary(tcur *dbank.TransactionSummary, tnew dbank.Transaction) error {
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

func (b *BankService) Transfer(fromAcct string, toAcct string, amount float64) (bool, error) {
	return true, nil
}
