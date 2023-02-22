package main

import (
	"log"

	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	grpcAdapter := mygrpc.NewAdapter(9090)
	grpcAdapter.Run()
}
