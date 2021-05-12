package auth

import (
	"context"
	"log"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		ur  userRepository
		jwt jwtInteractor
	}

	jwtInteractor interface {
		NewWithUserID(userID uuid.UUID) (uuid.UUID, string, error)
	}

	userRepository interface{}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(ji jwtInteractor, ur userRepository) *Service {
	if ur == nil {
		log.Fatalln("user repository is not set")
	}
	if ji == nil {
		log.Fatalln("jwt interactor is not set")
	}
	return &Service{jwt: ji, ur: ur}
}

// Login by email and password
func (s *Service) Login(ctx context.Context, email, password string) (string, error) {
	return "generated jwt string", nil
}
