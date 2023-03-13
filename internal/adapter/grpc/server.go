package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/timpamungkas/grpc-go-server/internal/port"
	"github.com/timpamungkas/grpc-proto/protogen/go/bank"
	"github.com/timpamungkas/grpc-proto/protogen/go/hello"
	rel "github.com/timpamungkas/grpc-proto/protogen/go/reliability"
	"google.golang.org/grpc"
)

type GrpcAdapter struct {
	helloService       port.HelloServicePort
	bankService        port.BankServicePort
	reliabilityService port.ReliabilityServicePort
	grpcPort           int
	server             *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
	rel.ReliabilityServiceServer
}

func NewGrpcAdapter(helloService port.HelloServicePort, bankService port.BankServicePort,
	reliabilityService port.ReliabilityServicePort, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		helloService:       helloService,
		bankService:        bankService,
		reliabilityService: reliabilityService,
		grpcPort:           grpcPort,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d : %v\n", a.grpcPort, err)
	}

	log.Printf("Server listening on %d\n", a.grpcPort)

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	hello.RegisterHelloServiceServer(grpcServer, a)
	bank.RegisterBankServiceServer(grpcServer, a)
	rel.RegisterReliabilityServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve grpc on %d : %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
