package application

var accounts map[string]int32

type BankService struct {
}

func init() {
	accounts = map[string]int32{
		"111": 5001,
		"222": 5002,
		"333": 5003,
	}
}

func (a *BankService) FindCurrentBalance(acct string) int32 {
	return accounts[acct]
}
