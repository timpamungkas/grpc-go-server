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

type BankService struct {
	db port.BankDatabasePort
}

func NewBankService(dbPort port.BankDatabasePort) *BankService {
	return &BankService{
		db: dbPort,
	}
}

func (b *BankService) FindCurrentBalance(acct string) float64 {
	bankAccount, err := b.db.GetBankAccountByAccountNumber(acct)

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
	exchangeRate, err := b.db.GetExchangeRateAtTimestamp(fromCur, toCur, ts)

	if err != nil {
		return 0
	}

	return float64(exchangeRate.Rate)
}

func (b *BankService) CreateTransaction(acct string, t dbank.Transaction) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	bankAccountOrm, err := b.db.GetBankAccountByAccountNumber(acct)

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

	savedUuid, err := b.db.CreateTransaction(bankAccountOrm, transactionOrm)

	if err != nil {
		return savedUuid, err
	}

	return savedUuid, nil
}

func (b *BankService) CalculateTransactionSummary(tcur *dbank.TransactionSummary, tnew dbank.Transaction) error {
	switch tnew.TransactionType {
	case dbank.TransactionStatusIn:
		tcur.SumIn += tnew.Amount
	case dbank.TransactionStatusOut:
		tcur.SumOut += tnew.Amount
	default:
		return fmt.Errorf("unknown transaction type : %v", tnew.TransactionType)
	}

	tcur.SumTotal = tcur.SumIn - tcur.SumOut

	return nil
}

func (b *BankService) Transfer(tt dbank.TransferTransaction) (uuid.UUID, bool, error) {
	now := time.Now()

	fromAccountOrm, err := b.db.GetBankAccountByAccountNumber(tt.FromAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer from account %v : %v", tt.FromAccountNumber, err)
		return uuid.Nil, false, err
	}

	toAccountOrm, err := b.db.GetBankAccountByAccountNumber(tt.ToAccountNumber)

	if err != nil {
		log.Printf("Can't find transfer to account %v : %v", tt.ToAccountNumber, err)
		return uuid.Nil, false, err
	}

	fromTransactionOrm := db.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		TransactionTimestamp: now,
		TransactionType:      dbank.TransactionStatusOut,
		AccountUuid:          fromAccountOrm.AccountUuid,
		Amount:               tt.Amount,
		Notes:                "Transfer out to " + tt.ToAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	toTransactionOrm := db.BankTransactionOrm{
		TransactionUuid:      uuid.New(),
		TransactionTimestamp: now,
		TransactionType:      dbank.TransactionStatusIn,
		AccountUuid:          toAccountOrm.AccountUuid,
		Amount:               tt.Amount,
		Notes:                "Transfer in from " + tt.FromAccountNumber,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	// create transfer request
	newTransferUuid := uuid.New()

	transferOrm := db.BankTransferOrm{
		TransferUuid:      newTransferUuid,
		FromAccountUuid:   fromAccountOrm.AccountUuid,
		ToAccountUuid:     toAccountOrm.AccountUuid,
		Currency:          tt.Currency,
		Amount:            tt.Amount,
		TransferTimestamp: now,
		TransferSuccess:   false,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if _, err := b.db.CreateTransfer(transferOrm); err != nil {
		log.Printf("Can't create transfer from %v to %v : %v",
			tt.FromAccountNumber, tt.ToAccountNumber, err)
		return uuid.Nil, false, err
	}

	if transferPairSuccess, err := b.db.CreateTransferTransactionPair(
		fromAccountOrm, toAccountOrm, fromTransactionOrm, toTransactionOrm); transferPairSuccess {
		b.db.UpdateTransferStatus(transferOrm, true)
		return newTransferUuid, true, nil
	} else {
		return newTransferUuid, false, err
	}
}
