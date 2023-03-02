package application

import "time"

type BankService struct {
}

func (a *BankService) FindCurrentBalance(acct string) int32 {
	return 999
}

func (a *BankService) FindExchangeRate(fromCur string, toCur string) float32 {
	now := time.Now()
	rate := int(1000) + now.Minute() + now.Second()

	return float32(rate)
}
