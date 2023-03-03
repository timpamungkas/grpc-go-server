package port

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) float64
	FindExchangeRate(fromCur string, toCur string) float64
}
