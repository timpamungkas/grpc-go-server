package port

import (
	"time"

	"github.com/google/uuid"
	db "github.com/timpamungkas/grpc-go-server/internal/adapter/database"
)

type DummyDatabasePort interface {
	Save(data *db.DummyOrm) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (db.DummyOrm, error)
}

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string, withTransactions bool,
		transactionFrom time.Time, transactionTo time.Time) (db.BankAccountOrm, error)
	CreateExchangeRate(r db.BankExchangeRateOrm) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (float64, error)
	CreateTransaction(t db.BankTransactionOrm) (uuid.UUID, error)
	UpdateCurrentBalance(acct db.BankAccountOrm, newBalance float64) error
}
