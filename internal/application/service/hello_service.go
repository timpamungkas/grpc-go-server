package application

import (
	"context"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a *Application) GenerateHello(ctx context.Context, name string) (string, error) {
	return "Hello " + name, nil
}
