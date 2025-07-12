package main

import (
	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/infra"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
)

func main() {
	infraProviders := infra.NewInfraProviders()
	infraProviders.UserRepository.CreateSessionsTable()
	infraProviders.UserRepository.CreateUserTable()
	up := uuid.NewProvider()
	en := bcrypt.NewProvider()
	appServices := app.NewServices(infraProviders.UserRepository, up, en)
	infraHTTPServer := infra.NewHTTPServer(appServices)
	infraHTTPServer.ListenAndServe(":8080")
}
