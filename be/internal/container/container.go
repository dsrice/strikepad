package container

import (
	"strikepad-backend/internal/auth"
	"strikepad-backend/internal/config"
	"strikepad-backend/internal/handler"
	"strikepad-backend/internal/repository"
	"strikepad-backend/internal/service"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	if err := container.Provide(config.NewDatabase); err != nil {
		panic(err)
	}
	if err := container.Provide(repository.NewUserRepository); err != nil {
		panic(err)
	}
	if err := container.Provide(repository.NewSessionRepository); err != nil {
		panic(err)
	}
	if err := container.Provide(auth.NewJWTService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewHealthService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewAPIService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewAuthService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewUserService); err != nil {
		panic(err)
	}
	if err := container.Provide(service.NewSessionService); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewHealthHandler); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewAPIHandler); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewAuthHandler); err != nil {
		panic(err)
	}
	if err := container.Provide(handler.NewUserHandler); err != nil {
		panic(err)
	}

	return container
}
