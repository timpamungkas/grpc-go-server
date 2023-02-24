package port

import "context"

type HelloServicePort interface {
	GenerateHello(ctx context.Context, name string) (string, error)
}
