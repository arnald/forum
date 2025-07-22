package app

import (
	"github.com/arnald/forum/internal/app/user/queries"
	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/bcrypt"
	"github.com/arnald/forum/internal/pkg/uuid"
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

func NewServices(repo user.Repository, up uuid.Provider, en bcrypt.Provider) Services {
	return Services{
		UserServices: UserServices{
			Queries: Queries{
				queries.NewUserRegisterHandler(repo, up, en),
			},
		},
	}
}
