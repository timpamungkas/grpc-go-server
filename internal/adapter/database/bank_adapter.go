package database

import (
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
