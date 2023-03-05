package database

import (
	"time"

	"github.com/google/uuid"
)

type BankAccount struct {
	AccountUuid    uuid.UUID `gorm:"primaryKey"`
	AccountNumber  string
	AccountName    string
	Currency       string
	CurrentBalance float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (BankAccount) TableName() string {
	return "bank_accounts"
}
