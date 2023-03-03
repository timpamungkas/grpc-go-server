package application

import "time"

var accounts map[string]float64

type BankService struct {
}

func init() {
	accounts = map[string]float64{
		"111": 5001,
		"222": 5002,
		"333": 5003,
	}
}

func (a *BankService) FindCurrentBalance(acct string) float64 {
	return accounts[acct]
}

func (a *BankService) FindExchangeRate(fromCur string, toCur string) float64 {
	now := time.Now()
	bal := 1000 + now.Minute() + now.Second()

	return float64(bal)
}
