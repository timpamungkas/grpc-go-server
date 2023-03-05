package port

import (
	"github.com/google/uuid"
	ddummy "github.com/timpamungkas/grpc-go-server/internal/application/domain/dummy"
)

type DummyDatabasePort interface {
	Save(data *ddummy.Dummy) (uuid.UUID, error)
	GetByUuid(uuid *uuid.UUID) (ddummy.Dummy, error)
}
