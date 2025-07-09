package app

import (
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/user"
)

type Queries struct {
	UserRegister queries.UserRegisterRequestHandler
}

type UserServices struct {
	Queries Queries
}

type Services struct {
	UserServices UserServices
}

func NewServices(repo user.Repository) Services {
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				queries.NewUserRegisterHandler(repo),
			},
		},
	}
}
