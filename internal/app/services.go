package app

type Queries struct{}

type UserServices struct {
	Queries Queries
}

type Services struct{}

func NewServices() Services {
	return Services{}
}
