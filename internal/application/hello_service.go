package application

import (
	"context"
)

func (a *Application) GenerateHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
