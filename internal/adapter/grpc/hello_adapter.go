package grpc

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
)

func (a *Adapter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	greet, err := a.helloService.GenerateHello(ctx, in.Name)

	if err != nil {
		return nil, err
	}

	return &pb.HelloResponse{
		Greet: greet,
	}, nil
}
