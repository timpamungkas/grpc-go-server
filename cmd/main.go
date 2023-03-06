package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	dbmigration "github.com/timpamungkas/grpc-go-server/db"
	mydb "github.com/timpamungkas/grpc-go-server/internal/adapter/database"
	mygrpc "github.com/timpamungkas/grpc-go-server/internal/adapter/grpc"
	app "github.com/timpamungkas/grpc-go-server/internal/application"
	dbank "github.com/timpamungkas/grpc-go-server/internal/application/domain/bank"
	ddummy "github.com/timpamungkas/grpc-go-server/internal/application/domain/dummy"
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

	go generateExchangeRates(databaseAdapter, "USD", "IDR", 5*time.Second)

	hs := new(app.HelloService)
	bs := app.NewBankService(databaseAdapter)
	grpcAdapter := mygrpc.NewGrpcAdapter(hs, bs, 9090)
	grpcAdapter.Run()
}

func runDummyOrm(da *mydb.DatabaseAdapter) {
	uuid, _ := da.Save(
		&ddummy.Dummy{
			UserName: "Tim " + time.Now().Format("15:04:05"),
		},
	)

	res, _ := da.GetByUuid(&uuid)
	log.Println("res : ", res)

}

func generateExchangeRates(da *mydb.DatabaseAdapter,
	fromCurrency string, toCurrency string, duration time.Duration) {
	ticker := time.NewTicker(duration)

	for range ticker.C {
		now := time.Now()
		validFrom := now.Truncate(time.Second)
		validTo := validFrom.Add(duration).Add(-1 * time.Millisecond)

		dummyRate := dbank.ExchangeRate{
			FromCurrency:       fromCurrency,
			ToCurrency:         toCurrency,
			ValidFromTimestamp: validFrom,
			ValidToTimestamp:   validTo,
			Rate:               2000 + float64(rand.Intn(300)),
		}

		da.CreateExchangeRate(dummyRate)
	}
}
