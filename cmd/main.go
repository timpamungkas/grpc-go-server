package main

import (
	"log"

	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
	svc "github.com/timpamungkas/grpc-go-server/internal/application/service"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	application := svc.NewApplication()
	grpcAdapter := mygrpc.NewAdapter(application, 9090)
	grpcAdapter.Run()
}
