package grpc

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
)

func (a *GrpcAdapter) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateHello(in.Name)

	return &hello.HelloResponse{
		Greet: greet,
	}, nil
}

func (a *GrpcAdapter) SayManyHellos(in *hello.HelloRequest, stream hello.HelloService_SayManyHellosServer) error {
	for i := 0; i < 10; i++ {
		greet := a.helloService.GenerateHello(in.Name)

		res := fmt.Sprintf("[%d] %s", i, greet)

		stream.Send(
			&hello.HelloResponse{
				Greet: res,
			},
		)

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (a *GrpcAdapter) SayHelloToEveryone(stream hello.HelloService_SayHelloToEveryoneServer) error {
	res := ""

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(
				&hello.HelloResponse{
					Greet: res,
				},
			)
		}

		if err != nil {
			log.Fatalf("Error while reading from client : %v", err)
		}

		greet := a.helloService.GenerateHello(req.Name)

		if err != nil {
			return err
		}

		res += greet + " "
	}
}

func (a *GrpcAdapter) SayHelloContinuous(stream hello.HelloService_SayHelloContinuousServer) error {
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error while reading from client : %v", err)
		}

		greet := a.helloService.GenerateHello(req.Name)

		if err != nil {
			return err
		}

		err = stream.Send(
			&hello.HelloResponse{
				Greet: greet,
			},
		)

		if err != nil {
			log.Fatalf("Error while sending response to client : %v", err)
		}
	}
}
