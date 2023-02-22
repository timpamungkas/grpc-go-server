package coursegrpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	port "github.com/timpamungkas/grpc-go-server/internal/port"
	"google.golang.org/grpc"
)

type Adapter struct {
	service  port.ServicePort
	grpcPort int
	server   *grpc.Server
	pb.HelloServiceServer
}

func NewAdapter(service port.ServicePort, grpcPort int) *Adapter {
	return &Adapter{
		service:  service,
		grpcPort: grpcPort,
	}
}

func (a *Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d : %v\n", a.grpcPort, err)
	}

	log.Printf("Server listening on %d\n", a.grpcPort)

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	pb.RegisterHelloServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve grpc on %d : %v\n", a.grpcPort, err)
	}
}

func (a *Adapter) Stop() {
	a.server.Stop()
}
