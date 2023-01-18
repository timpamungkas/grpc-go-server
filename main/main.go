package main

import (
	"log"
	"net"

	pb "github.com/timpamungkas/course-grpc-proto/protogen/go/hello"
	cg "github.com/timpamungkas/grpc-go-server/coursegrpc"
	"google.golang.org/grpc"
)

const serverAddr string = ":9090"

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	lis, err := net.Listen("tcp", serverAddr)

	if err != nil {
		log.Fatalf("Failed to listen on %v : %v\n", serverAddr, err)
	}

	log.Printf("Server listening on %v\n", serverAddr)

	server := grpc.NewServer()

	pb.RegisterHelloServiceServer(server, &cg.GrpcServer{})

	if err = server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpc on %v : %v\n", serverAddr, err)
	}
}
