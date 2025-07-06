package infra

import (
	"github.com/arnald/forum/internal/infra/http"
)

type Services struct {
	// UserRepository user.Repository
	Server *http.Server
}

func NewInfraProviders() Services {
	return Services{
		// UserRepository: sqlite.NewRepo(),
	}
}

func NewHTTPServer() *http.Server {
	return http.NewServer()
}
