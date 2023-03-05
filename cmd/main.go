package main

import (
	"database/sql"
	"log"
	"time"

	dbmigration "github.com/timpamungkas/grpc-go-server/db"
	mydb "github.com/timpamungkas/grpc-go-server/internal/adapter/database"
	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
	app "github.com/timpamungkas/grpc-go-server/internal/application"
	"github.com/timpamungkas/grpc-go-server/internal/application/domain/dummy"
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

	// runDummyOrm(databaseAdapter)

	hs := new(app.HelloService)
	bs := app.NewBankService(databaseAdapter)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()
}

func runDummyOrm(da *mydb.DatabaseAdapter) {
	uuid, _ := da.Save(
		&dummy.Dummy{
			UserName: "Tim " + time.Now().Format("15:04:05"),
		},
	)

	res, _ := da.GetByUuid(&uuid)
	log.Println("res : ", res)

}
