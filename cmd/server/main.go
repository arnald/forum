package main

import (
	"github.com/arnald/forum/internal/infra"
)

func main() {
	// infraProviders := infra.NewInfraProviders()
	// appServices := app.NewServices()
	infraHTTPServer := infra.NewHTTPServer()
	infraHTTPServer.ListenAndServe()
}
