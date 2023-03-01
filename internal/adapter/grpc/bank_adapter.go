package grpc

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/bank"
)

func (a *GrpcAdapter) GetCurrentBalance(
	ctx context.Context, in *pb.CurrentBalanceRequest) (*pb.CurrentBalanceResponse, error) {
	bal := a.bankService.FindCurrentBalance(in.AccountNumber)

	return &pb.CurrentBalanceResponse{
		Amount: bal,
	}, nil
}
