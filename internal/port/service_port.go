package port

type HelloServicePort interface {
	GenerateHello(name string) string
}

type BankServicePort interface {
	FindCurrentBalance(acct string) int32
	FindExchangeRate(fromCur string, toCur string) float32
}
