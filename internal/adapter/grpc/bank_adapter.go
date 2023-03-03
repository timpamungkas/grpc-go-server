package grpc

import (
	"context"
	"time"

	"github.com/timpamungkas/grpc-proto/protogen/go/bank"
	"google.golang.org/genproto/googleapis/type/date"
)

func (a *GrpcAdapter) GetCurrentBalance(
	ctx context.Context, in *bank.CurrentBalanceRequest) (*bank.CurrentBalanceResponse, error) {
	now := time.Now()
	bal := a.bankService.FindCurrentBalance(in.AccountNumber)

	return &bank.CurrentBalanceResponse{
		Amount: bal,
		CurrentDate: &date.Date{
			Year:  int32(now.Year()),
			Month: int32(now.Month()),
			Day:   int32(now.Day()),
		},
	}, nil
}

func (a *GrpcAdapter) FetchExchangeRate(in *bank.ExchangeRateRequest,
	stream bank.BankService_FetchExchangeRateServer) error {
	for {
		rate := a.bankService.FindExchangeRate(in.FromCurrency, in.ToCurrency)

		stream.Send(
			&bank.ExchangeRateResponse{
				FromCurrency: in.FromCurrency,
				ToCurrency:   in.ToCurrency,
				Rate:         rate,
				Timestamp:    time.Now().Format(time.RFC3339),
			},
		)

		time.Sleep(1 * time.Second)
	}
}
