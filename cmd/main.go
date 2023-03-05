package main

import (
	"database/sql"
	"log"

	dbmigration "github.com/timpamungkas/grpc-go-server/db"
	mydb "github.com/timpamungkas/grpc-go-server/internal/adapter/database"
	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
	"github.com/timpamungkas/grpc-go-server/internal/application"
	app "github.com/timpamungkas/grpc-go-server/internal/application"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(logWriter{})

	sqlDB, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/grpc?sslmode=disable")

	if err != nil {
		log.Fatalf("Can't open connection to database : %v\n", err)
	}

	dbmigration.Migrate(sqlDB)

	databaseAdapter, err := mydb.NewDatabaseAdapter(sqlDB)

	if err != nil {
		log.Fatalf("Can't create database adapter : %v\n", err)
	}

	// databaseAdapter.Save(
	// 	&dummy.Dummy{
	// 		UserName: "Zoe",
	// 	},
	// )

	hs := new(app.HelloService)
	bs := application.NewBankService(databaseAdapter)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()
}
