package database

import (
	"time"

	"github.com/google/uuid"
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
)

func (a *DatabaseAdapter) GetBankAccountByAccountNumber(acct string, withTransactions bool) (dbank.Account, error) {
	var bankAccountOrm BankAccountOrm
	var res dbank.Account

	err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error

	res = dbank.Account{
		AccountNumber:  bankAccountOrm.AccountNumber,
		AccountName:    bankAccountOrm.AccountName,
		Currency:       bankAccountOrm.Currency,
		CurrentBalance: bankAccountOrm.CurrentBalance,
	}

	return res, err
}

func (a *DatabaseAdapter) CreateExchangeRate(r dbank.ExchangeRate) (uuid.UUID, error) {
	newUuid := uuid.New()
	now := time.Now()

	dummyRate := BankExchangeRateOrm{
		ExchangeRateUuid:   newUuid,
		FromCurrency:       r.FromCurrency,
		ToCurrency:         r.ToCurrency,
		ValidFromTimestamp: r.ValidFromTimestamp,
		ValidToTimestamp:   r.ValidToTimestamp,
		Rate:               r.Rate,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := a.db.Create(dummyRate).Error; err != nil {
		return uuid.Nil, err
	}

	return newUuid, nil
}

func (a *DatabaseAdapter) GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (float64, error) {
	var exchangeRateOrm BankExchangeRateOrm

	err := a.db.First(&exchangeRateOrm, "from_currency = ? "+
		" AND to_currency = ? "+
		" AND (? BETWEEN valid_from_timestamp AND valid_to_timestamp)",
		fromCur, toCur, ts).Error

	return exchangeRateOrm.Rate, err
}
