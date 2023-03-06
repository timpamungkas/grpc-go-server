package application

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	db "github.com/timpamungkas/grpc-go-server/internal/adapter/database"
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
	bankAccount, err := b.db.GetBankAccountByAccountNumber(acct, false, time.Now(), time.Now())

	if err != nil {
		log.Printf("Error on FindCurrentBalance : %v\n", err)
	}

	return bankAccount.CurrentBalance
}

func (b *BankService) CreateExchangeRate(r dbank.ExchangeRate) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	exchangeRateOrm := db.BankExchangeRateOrm{
		ExchangeRateUuid:   newUuid,
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		Rate:               r.Rate,
		ValidFromTimestamp: r.ValidFromTimestamp,
		ValidToTimestamp:   r.ValidToTimestamp,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	return b.db.CreateExchangeRate(exchangeRateOrm)
}

func (b *BankService) FindExchangeRate(fromCur string, toCur string, ts time.Time) float64 {
	rate, err := b.db.GetExchangeRateAtTimestamp(fromCur, toCur, ts)

	if err != nil {
		return 0
	}

	return float64(rate)
}

func (b *BankService) CreateTransaction(acct string, t dbank.Transaction) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	bankAccountOrm, err := b.db.GetBankAccountByAccountNumber(acct, false, time.Now(), time.Now())

	if err != nil {
		log.Printf("Can't create transaction for %v : %v", acct, err)
		return uuid.Nil, err
	}

	transactionOrm := db.BankTransactionOrm{
		TransactionUuid:      newUuid,
		AccountUuid:          bankAccountOrm.AccountUuid,
		TransactionTimestamp: now,
		Amount:               t.Amount,
		TransactionType:      t.TransactionType,
		Notes:                t.Notes,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	savedUuid, err := b.db.CreateTransaction(transactionOrm)

	if err != nil {
		return savedUuid, err
	}

	// recalculate current balance
	newAmount := t.Amount

	if t.TransactionType == dbank.Out {
		newAmount = -1 * t.Amount
	}

	newAccountBalance := bankAccountOrm.CurrentBalance + newAmount

	b.db.UpdateCurrentBalance(bankAccountOrm, newAccountBalance)

	return savedUuid, nil
}

func (b *BankService) CalculateTransactionSummary(tcur *dbank.TransactionSummary, tnew dbank.Transaction) error {
	switch tnew.TransactionType {
	case dbank.In:
		tcur.SumIn += tnew.Amount
	case dbank.Out:
		tcur.SumOut += tnew.Amount
	default:
		return fmt.Errorf("unknown transaction type : %v", tnew.TransactionType)
	}

	tcur.SumTotal = tcur.SumIn - tcur.SumOut

	return nil
}

func (b *BankService) Transfer(fromAcct string, toAcct string, amount float64) (bool, error) {
	return true, nil
}
