package shows

import (
	"context"
	"log"

	"github.com/SatorNetwork/sator-api/svc/shows/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		sr showsRepository
	}

	// Show struct
	Show struct {
		ID            uuid.UUID `json:"id"`
		Title         string    `json:"title"`
		Cover         string    `json:"cover"`
		HasNewEpisode bool      `json:"has_new_episode"`
	}

	showsRepository interface {
		GetShows(ctx context.Context, arg repository.GetShowsParams) ([]repository.Show, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(sr showsRepository) *Service {
	if sr == nil {
		log.Fatalln("shows repository is not set")
	}
	return &Service{sr: sr}
}

// GetShows returns shows.
func (s *Service) GetShows(ctx context.Context, page int) (interface{}, error) {
	shows, err := s.sr.GetShows(ctx, repository.GetShowsParams{Page: page})
	if err != nil {
		return []repository.Show{}, err
	}

	return shows, nil
}
