package firebase

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	firebase_repository "github.com/SatorNetwork/sator-api/svc/firebase/repository"
)

type (
	Service struct {
		fr firebaseRepository
	}

	firebaseRepository interface {
		GetRegistrationToken(
			ctx context.Context,
			arg firebase_repository.GetRegistrationTokenParams,
		) (firebase_repository.FirebaseRegistrationToken, error)
		UpsertRegistrationToken(
			ctx context.Context,
			arg firebase_repository.UpsertRegistrationTokenParams,
		) error
	}

	Empty struct{}

	RegisterTokenRequest struct {
		DeviceId string `json:"device_id"`
		Token    string `json:"token"`
	}
)

func NewService(
	fr firebaseRepository,
) *Service {
	s := &Service{
		fr: fr,
	}

	return s
}

func (s *Service) RegisterToken(ctx context.Context, userID uuid.UUID, req *RegisterTokenRequest) (*Empty, error) {
	err := s.fr.UpsertRegistrationToken(ctx, firebase_repository.UpsertRegistrationTokenParams{
		UserID:            userID,
		DeviceID:          req.DeviceId,
		RegistrationToken: req.Token,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't register token")
	}

	return &Empty{}, nil
}
