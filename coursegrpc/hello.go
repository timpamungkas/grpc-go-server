package coursegrpc

import (
	"context"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
)

func (s *GrpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Greet: "Hello " + in.Name,
	}, nil
}
