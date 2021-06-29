package shows

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/svc/shows/repository"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		sr  showsRepository
		chc challengesClient
	}

	// Show struct
	// Fields were rearranged to optimize memory usage.
	Show struct {
		ID            uuid.UUID `json:"id"`
		Title         string    `json:"title"`
		Cover         string    `json:"cover"`
		HasNewEpisode bool      `json:"has_new_episode"`
	}

	Episode struct {
		ID            uuid.UUID `json:"id"`
		ShowID        uuid.UUID `json:"show_id"`
		EpisodeNumber int32     `json:"episode_number"`
		Cover         string    `json:"cover"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		ReleaseDate   string    `json:"release_date"`
	}

	showsRepository interface {
		GetShows(ctx context.Context, arg repository.GetShowsParams) ([]repository.Show, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.Show, error)
		GetShowsByCategory(ctx context.Context, arg repository.GetShowsByCategoryParams) ([]repository.Show, error)
		GetEpisodeByID(ctx context.Context, id uuid.UUID) (repository.Episode, error)
		GetEpisodesByShowID(ctx context.Context, arg repository.GetEpisodesByShowIDParams) ([]repository.Episode, error)
	}

	// Challenges service client
	challengesClient interface {
		GetListByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(sr showsRepository, chc challengesClient) *Service {
	if sr == nil {
		log.Fatalln("shows repository is not set")
	}
	if chc == nil {
		log.Fatalln("challenges client is not set")
	}

	return &Service{sr: sr, chc: chc}
}

// GetShows returns shows.
func (s *Service) GetShows(ctx context.Context, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetShows(ctx, repository.GetShowsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}
	return castToListShow(shows), nil
}

// GetShowChallenges returns challenges by show id.
func (s *Service) GetShowChallenges(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	challenges, err := s.chc.GetListByShowID(ctx, showID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not get challenges list by show id: %w", err)
	}

	return challenges, nil
}

// Cast repository.Show to service Show structure
func castToListShow(source []repository.Show) []Show {
	result := make([]Show, 0, len(source))
	for _, s := range source {
		result = append(result, Show{
			ID:            s.ID,
			Title:         s.Title,
			Cover:         s.Cover,
			HasNewEpisode: s.HasNewEpisode,
		})
	}
	return result
}

// GetShowByID returns show with provided id.
func (s *Service) GetShowByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	show, err := s.sr.GetShowByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get show with id=%s: %w", id, err)
	}
	return castToShow(show), nil
}

// Cast repository.Show to service Show structure
func castToShow(source repository.Show) Show {
	return Show{
		ID:            source.ID,
		Title:         source.Title,
		Cover:         source.Cover,
		HasNewEpisode: source.HasNewEpisode,
	}
}

// GetShowsByCategory returns show by provided category.
func (s *Service) GetShowsByCategory(ctx context.Context, category string, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetShowsByCategory(ctx, repository.GetShowsByCategoryParams{
		Category: sql.NullString{String: category, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}
	return castToListShow(shows), nil
}

// GetEpisodesByShowID returns episodes by show id.
func (s *Service) GetEpisodesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	episodes, err := s.sr.GetEpisodesByShowID(ctx, repository.GetEpisodesByShowIDParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get episodes list: %w", err)
	}
	return castToListEpisodes(episodes), nil
}

// Cast repository.Episode to service Episode structure
func castToListEpisodes(source []repository.Episode) []Episode {
	result := make([]Episode, 0, len(source))
	for _, s := range source {
		result = append(result, Episode{
			ID:            s.ID,
			ShowID:        s.ShowID,
			EpisodeNumber: s.EpisodeNumber,
			Cover:         s.Cover.String,
			Title:         s.Title,
			Description:   s.Description.String,
			ReleaseDate:   s.ReleaseDate.Time.String(),
		})
	}
	return result
}

// GetEpisodeByID returns episode with provided id.
func (s *Service) GetEpisodeByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	episode, err := s.sr.GetEpisodeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get episode with id=%s: %w", id, err)
	}
	return castToEpisode(episode), nil
}

// Cast repository.Episode to service Episode structure
func castToEpisode(source repository.Episode) Episode {
	return Episode{
		ID:            source.ID,
		ShowID:        source.ShowID,
		EpisodeNumber: source.EpisodeNumber,
		Cover:         source.Cover.String,
		Title:         source.Title,
		Description:   source.Description.String,
		ReleaseDate:   source.ReleaseDate.Time.String(),
	}
}
