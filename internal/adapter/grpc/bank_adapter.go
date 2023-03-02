package grpc

import (
	"context"
	"time"

	"github.com/timpamungkas/course-grpc-proto/protogen/go/bank"
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
