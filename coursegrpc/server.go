package coursegrpc

import (
	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
)

type GrpcServer struct {
	pb.HelloServiceServer
}
