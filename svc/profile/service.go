package profile

import (
	"context"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/profile/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		pr profileRepository
	}

	profileRepository interface {
		CreateProfile(ctx context.Context, arg repository.CreateProfileParams) (repository.Profile, error)
		GetProfileByUserID(ctx context.Context, userID uuid.UUID) (repository.Profile, error)
	}

	// Profile struct
	Profile struct {
		UserID    string `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
)

// NewService is a factory function, returns a new instance of the Service interface implementation
func NewService(pr profileRepository) *Service {
	if pr == nil {
		log.Fatalln("profile repository is not set")
	}
	return &Service{pr: pr}
}

// GetProfileByUserID returns user profile by user id
// If it doesn't exist, the method creates missed record in database and returns it
func (s *Service) GetProfileByUserID(ctx context.Context, userID uuid.UUID) (interface{}, error) {
	profile, err := s.pr.GetProfileByUserID(ctx, userID)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, fmt.Errorf("could not get user profile: %w", err)
		}

		profile, err = s.pr.CreateProfile(ctx, repository.CreateProfileParams{
			UserID: userID,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get user profile: %w", err)
		}
	}

	return castToProfile(profile), nil
}

func castToProfile(p repository.Profile) *Profile {
	return &Profile{
		UserID:    p.UserID.String(),
		FirstName: p.FirstName.String,
		LastName:  p.LastName.String,
	}
}
