package main

import (
	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/infra"
	"github.com/arnald/forum/internal/pkg/uuid"
)

func main() {
	infraProviders := infra.NewInfraProviders()
	up := uuid.NewProvider()
	appServices := app.NewServices(infraProviders.UserRepository, up)
	infraHTTPServer := infra.NewHTTPServer(appServices)
	infraHTTPServer.ListenAndServe(":8080")
}
