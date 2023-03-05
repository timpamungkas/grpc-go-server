package database

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseAdapter struct {
	db *gorm.DB
}

func NewDatabaseAdapter(conn *sql.DB) (*DatabaseAdapter, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: conn,
	}), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("can't connect database : %v", err)
	}

	err = db.AutoMigrate(&DummyOrm{})

	if err != nil {
		return nil, fmt.Errorf("database migration error: %v", err)
	}

	return &DatabaseAdapter{db: db}, nil
}
