package database

import (
	"log"
	"time"

	"github.com/google/uuid"
	domain "github.com/timpamungkas/grpc-go-server/internal/application/domain/dummy"
)

func (a *DatabaseAdapter) Save(data *domain.Dummy) (uuid.UUID, error) {
	now := time.Now()
	userId := data.UserId

	if data.UserId == uuid.Nil {
		userId = uuid.New()
	}

	dummyData := DummyOrm{
		UserId:    userId,
		UserName:  data.UserName,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := a.db.Create(dummyData).Error

	if err != nil {
		log.Printf("Can't create data : %v", err)
	}

	return userId, nil
}
