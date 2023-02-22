package coursegrpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	"google.golang.org/grpc"
)

type Adapter struct {
	port   int
	server *grpc.Server
	pb.HelloServiceServer
}

func NewAdapter(port int) *Adapter {
	return &Adapter{
		port: port,
	}
}

func (a Adapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d : %v\n", a.port, err)
	}

	log.Printf("Server listening on %d\n", a.port)

	grpcServer := grpc.NewServer()
	a.server = grpcServer

	pb.RegisterHelloServiceServer(grpcServer, a)

	if err = grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve grpc on %d : %v\n", a.port, err)
	}
}

func (a Adapter) Stop() {
	a.server.Stop()
}
