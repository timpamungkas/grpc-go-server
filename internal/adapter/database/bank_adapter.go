package database

import (
	"log"
	"time"

	"github.com/google/uuid"
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
)

func (a *DatabaseAdapter) GetBankAccountByAccountNumber(acct string) (BankAccountOrm, error) {
	var bankAccountOrm BankAccountOrm

	if err := a.db.First(&bankAccountOrm, "account_number = ?", acct).Error; err != nil {
		log.Printf("Can't find bank account %v : %v", acct, err)
		return bankAccountOrm, err
	}

	return bankAccountOrm, nil
}

func (a *DatabaseAdapter) CreateExchangeRate(r BankExchangeRateOrm) (uuid.UUID, error) {
	if err := a.db.Create(r).Error; err != nil {
		return uuid.Nil, err
	}

	return r.ExchangeRateUuid, nil
}

func (a *DatabaseAdapter) GetExchangeRateAtTimestamp(fromCur string, toCur string, ts time.Time) (BankExchangeRateOrm, error) {
	var exchangeRateOrm BankExchangeRateOrm

	err := a.db.First(&exchangeRateOrm, "from_currency = ? "+
		" AND to_currency = ? "+
		" AND (? BETWEEN valid_from_timestamp AND valid_to_timestamp)",
		fromCur, toCur, ts).Error

	return exchangeRateOrm, err
}

func (a *DatabaseAdapter) CreateTransaction(acct BankAccountOrm, t BankTransactionOrm) (uuid.UUID, error) {
	tx := a.db.Begin()

	if err := tx.Create(t).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	// recalculate current balance
	newAmount := t.Amount

	if t.TransactionType == dbank.TransactionStatusOut {
		newAmount = -1 * t.Amount
	}

	newAccountBalance := acct.CurrentBalance + newAmount

	if err := tx.Model(&acct).Updates(
		map[string]interface{}{
			"current_balance": newAccountBalance,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()

	return t.TransactionUuid, nil
}

func (a *DatabaseAdapter) CreateTransfer(transfer BankTransferOrm) (uuid.UUID, error) {
	if err := a.db.Create(transfer).Error; err != nil {
		return uuid.Nil, err
	}

	return transfer.TransferUuid, nil
}

func (a *DatabaseAdapter) CreateTransferTransactionPair(
	fromAccountOrm BankAccountOrm,
	toAccountOrm BankAccountOrm,
	fromTransactionOrm BankTransactionOrm,
	toTransactionOrm BankTransactionOrm) (bool, error) {
	tx := a.db.Begin()

	if err := tx.Create(fromTransactionOrm).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Create(toTransactionOrm).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// recalculate current balance (fromAccount)
	fromAccountNewBalance := fromAccountOrm.CurrentBalance - fromTransactionOrm.Amount

	if err := tx.Model(&fromAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": fromAccountNewBalance,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// recalculate current balance (toAccount)
	toAccountNewBalance := toAccountOrm.CurrentBalance + toTransactionOrm.Amount

	if err := tx.Model(&toAccountOrm).Updates(
		map[string]interface{}{
			"current_balance": toAccountNewBalance,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()

	return true, nil
}

func (a *DatabaseAdapter) UpdateTransferStatus(transfer BankTransferOrm, status bool) error {
	if err := a.db.Model(&transfer).Updates(
		map[string]interface{}{
			"transfer_status": status,
			"updated_at":      time.Now(),
		},
	).Error; err != nil {
		return err
	}

	return nil
}
