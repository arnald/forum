package queries

import (
	"log"
	"time"

	"github.com/arnald/forum/internal/domain/user"
	"github.com/arnald/forum/internal/pkg/uuid"
)

type CreateSessionRequest struct {
	UserID    string
	Expiry    time.Time
	IPAddress string
}

type CreateSessionRequestHandler interface {
	Handle(req CreateSessionRequest) (*user.Session, error)
}

type createSessionRequestHandler struct {
	repo         user.Repository
	uuidProvider uuid.Provider
}

func NewCreateSessionRequestHandler(repo user.Repository, uuidProvider uuid.Provider) createSessionRequestHandler {
	return createSessionRequestHandler{
		repo:         repo,
		uuidProvider: uuidProvider,
	}
}

func (h createSessionRequestHandler) Handle(req CreateSessionRequest) (*user.Session, error) {
	session := &user.Session{
		Token:     []byte(h.uuidProvider.NewUUID()),
		UserID:    req.UserID,
		Expiry:    req.Expiry,
		IPAddress: req.IPAddress,
	}

	err := h.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}

	log.Println("Session created successfully")

	return session, nil
}
