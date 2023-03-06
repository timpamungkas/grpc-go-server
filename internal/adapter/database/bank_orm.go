package database

import (
	"time"

	"github.com/google/uuid"
)

type BankAccountOrm struct {
	AccountUuid    uuid.UUID `gorm:"primaryKey"`
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (BankAccountOrm) TableName() string {
	return "bank_accounts"
}

type BankExchangeRateOrm struct {
	ExchangeRateUuid   uuid.UUID `gorm:"primaryKey"`
	FromCurrency       string
	ToCurrency         string
	Rate               float64
	ValidFromTimestamp time.Time
	ValidToTimestamp   time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (BankExchangeRateOrm) TableName() string {
	return "bank_exchange_rates"
}
