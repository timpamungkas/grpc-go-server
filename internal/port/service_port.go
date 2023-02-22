package ports

import "context"

type ServicePort interface {
	GenerateHello(ctx context.Context, name string) (string, error)
}
