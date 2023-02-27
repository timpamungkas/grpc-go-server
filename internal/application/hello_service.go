package application

func (a *Application) GenerateHello(name string) (string, error) {
	return "Hello " + name, nil
}
