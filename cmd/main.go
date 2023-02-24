package main

import (
	"log"

	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
	app "github.com/timpamungkas/grpc-go-server/internal/application"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	application := app.NewApplication()
	grpcAdapter := mygrpc.NewGrpcAdapter(application, 9090)
	grpcAdapter.Run()
}
