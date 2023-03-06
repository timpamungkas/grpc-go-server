package database

import (
	"log"
	"time"

	"github.com/google/uuid"
)

func (a *DatabaseAdapter) GetBankAccountByAccountNumber(
	acct string, withTransactions bool, transactionFrom time.Time,
	transactionTo time.Time) (BankAccountOrm, error) {
	var bankAccountOrm BankAccountOrm

	if err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error; err != nil {
		log.Printf("Can't find bank account %v : %v", acct, err)
		return bankAccountOrm, err
	}

	if withTransactions {
		var txOrm []BankTransactionOrm

		a.db.Order("transaction_timestamp DESC").
			Find(&txOrm, "account_uuid = ? "+
				" AND transaction_timestamp BETWEEN ? AND ?",
				bankAccountOrm.AccountUuid, transactionFrom, transactionTo)

		bankAccountOrm.Transactions = append(bankAccountOrm.Transactions, txOrm...)
	}

	return bankAccountOrm, nil
}

func (a *DatabaseAdapter) CreateExchangeRate(r BankExchangeRateOrm) (uuid.UUID, error) {
	if err := a.db.Create(r).Error; err != nil {
		return uuid.Nil, err
	}

	return r.ExchangeRateUuid, nil
}

func (a *DatabaseAdapter) GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (float64, error) {
	var exchangeRateOrm BankExchangeRateOrm

	err := a.db.First(&exchangeRateOrm, "from_currency = ? "+
		" AND to_currency = ? "+
		" AND (? BETWEEN valid_from_timestamp AND valid_to_timestamp)",
		fromCur, toCur, ts).Error

	return exchangeRateOrm.Rate, err
}

func (a *DatabaseAdapter) CreateTransaction(t BankTransactionOrm) (uuid.UUID, error) {
	if err := a.db.Create(t).Error; err != nil {
		return uuid.Nil, err
	}

	return t.TransactionUuid, nil
}

func (a *DatabaseAdapter) UpdateCurrentBalance(acct BankAccountOrm, newBalance float64) error {
	return a.db.Model(&acct).Updates(
		map[string]interface{}{
			"current_balance": newBalance,
			"updated_at":      time.Now(),
		},
	).Error
}
