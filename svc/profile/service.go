package profile

import (
	"context"
	"database/sql"
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
		UpdateAvatar(ctx context.Context, arg repository.UpdateAvatarParams) error
	}

	// Profile struct
	Profile struct {
		UserID    string `json:"id"`
		UserName  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
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
func (s *Service) GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (interface{}, error) {
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

	return castToProfile(profile, username), nil
}

func castToProfile(p repository.Profile, username string) *Profile {
	profile := &Profile{
		UserID:    p.UserID.String(),
		UserName:  username,
		FirstName: p.FirstName.String,
		LastName:  p.LastName.String,
		Avatar:    p.Avatar.String,
	}

	return profile
}

// UpdateAvatar updates users avatar.
func (s *Service) UpdateAvatar(ctx context.Context, uid uuid.UUID, avatar string) error {
	_, err := s.pr.GetProfileByUserID(ctx, uid)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return fmt.Errorf("could not get user profile: %w", err)
		}

		_, err = s.pr.CreateProfile(ctx, repository.CreateProfileParams{
			UserID: uid,
		})
		if err != nil {
			return fmt.Errorf("could not get user profile: %w", err)
		}
	}

	if err := s.pr.UpdateAvatar(ctx, repository.UpdateAvatarParams{
		Avatar: sql.NullString{
			String: avatar,
			Valid:  len(avatar) > 0,
		},
		UserID: uid,
	}); err != nil {
		return err
	}

	return nil
}
