package main

import (
	"github.com/arnald/forum/internal/app"
	"github.com/arnald/forum/internal/infra"
)

func main() {
	infraProviders := infra.NewInfraProviders()
	appServices := app.NewServices(infraProviders.UserRepository)
	infraHTTPServer := infra.NewHTTPServer(appServices)
	infraHTTPServer.ListenAndServe()
}
