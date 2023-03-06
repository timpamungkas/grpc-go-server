package port

import (
	"time"

	"github.com/google/uuid"
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
	ddummy "github.com/timpamungkas/grpc-go-server/internal/application/domain/dummy"
)

type DummyDatabasePort interface {
	Save(data *ddummy.Dummy) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (ddummy.Dummy, error)
}

type BankDatabasePort interface {
	GetBankAccountByAccountNumber(acct string, withTransactions bool) (dbank.Account, error)
	CreateExchangeRate(r dbank.ExchangeRate) (uuid.UUID, error)
	GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (float64, error)
}
