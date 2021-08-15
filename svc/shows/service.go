package shows

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/internal/utils"
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
		Category      string    `json:"category"`
		Description   string    `json:"description"`
	}

	Season struct {
		ID           uuid.UUID `json:"id"`
		Title        string    `json:"title"`
		SeasonNumber int       `json:"season_number"`
		Episodes     []Episode `json:"episodes"`
	}

	Episode struct {
		ID            uuid.UUID `json:"id"`
		ShowID        uuid.UUID `json:"show_id"`
		EpisodeNumber int32     `json:"episode_number"`
		Cover         string    `json:"cover"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		ReleaseDate   string    `json:"release_date"`
		ChallengeID   uuid.UUID `json:"challenge_id"`
	}

	showsRepository interface {
		AddShow(ctx context.Context, arg repository.AddShowParams) error
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetShows(ctx context.Context, arg repository.GetShowsParams) ([]repository.Show, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.Show, error)
		GetShowsByCategory(ctx context.Context, arg repository.GetShowsByCategoryParams) ([]repository.Show, error)
		UpdateShow(ctx context.Context, arg repository.UpdateShowParams) error

		AddEpisode(ctx context.Context, arg repository.AddEpisodeParams) error
		GetEpisodeByID(ctx context.Context, arg repository.GetEpisodeByIDParams) (repository.Episode, error)
		GetEpisodesByShowID(ctx context.Context, arg repository.GetEpisodesByShowIDParams) ([]repository.Episode, error)
		DeleteEpisodeByID(ctx context.Context, arg repository.DeleteEpisodeByIDParams) error
		UpdateEpisode(ctx context.Context, arg repository.UpdateEpisodeParams) error

		GetSeasonsByShowID(ctx context.Context, arg repository.GetSeasonsByShowIDParams) ([]repository.Season, error)
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
			Category:      s.Category.String,
			Description:   s.Description.String,
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
		Category:      source.Category.String,
		Description:   source.Description.String,
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
	seasons, err := s.sr.GetSeasonsByShowID(ctx, repository.GetSeasonsByShowIDParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get seasons list: %w", err)
	}

	episodes, err := s.sr.GetEpisodesByShowID(ctx, repository.GetEpisodesByShowIDParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get episodes list: %w", err)
	}

	episodesPerSeasons := make(map[string][]Episode)
	for _, e := range episodes {
		if _, ok := episodesPerSeasons[e.SeasonID.String()]; ok {
			episodesPerSeasons[e.SeasonID.String()] = append(episodesPerSeasons[e.SeasonID.String()], castToEpisode(e))
		} else {
			episodesPerSeasons[e.SeasonID.String()] = []Episode{castToEpisode(e)}
		}
	}

	return castToListSeasons(seasons, episodesPerSeasons), nil
}

// Cast repository.Season to service Season structure
func castToListSeasons(source []repository.Season, episodes map[string][]Episode) []Season {
	result := make([]Season, 0, len(source))
	for _, s := range source {
		result = append(result, Season{
			ID:           s.ID,
			SeasonNumber: int(s.SeasonNumber),
			Title:        fmt.Sprintf("Season %d", s.SeasonNumber),
			Episodes:     episodes[s.ID.String()],
		})
	}
	return result
}

// Cast repository.Episode to service Episode structure
// func castToListEpisodes(source []repository.Episode) []Episode {
// 	result := make([]Episode, 0, len(source))
// 	for _, s := range source {
// 		result = append(result, Episode{
// 			ID:            s.ID,
// 			ShowID:        s.ShowID,
// 			EpisodeNumber: s.EpisodeNumber,
// 			Cover:         s.Cover.String,
// 			Title:         s.Title,
// 			Description:   s.Description.String,
// 			ReleaseDate:   s.ReleaseDate.Time.String(),
// 			ChallengeID:   s.ChallengeID,
// 		})
// 	}
// 	return result
// }

// GetEpisodeByID returns episode with provided id.
func (s *Service) GetEpisodeByID(ctx context.Context, showID, episodeID uuid.UUID) (interface{}, error) {
	episode, err := s.sr.GetEpisodeByID(ctx, repository.GetEpisodeByIDParams{
		ID:     episodeID,
		ShowID: showID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get episode with id=%s: %w", episodeID, err)
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
		ChallengeID:   source.ChallengeID,
	}
}

// AddShow ..
func (s *Service) AddShow(ctx context.Context, sh Show) error {
	if err := s.sr.AddShow(ctx, repository.AddShowParams{
		Title:         sh.Title,
		Cover:         sh.Cover,
		HasNewEpisode: sh.HasNewEpisode,
		Category: sql.NullString{
			String: sh.Category,
			Valid:  true,
		},
		Description: sql.NullString{
			String: sh.Description,
			Valid:  len(sh.Description) > 0,
		},
	}); err != nil {
		return fmt.Errorf("could not add show with title=%s: %w", sh.Title, err)
	}

	return nil
}

// UpdateShow ..
func (s *Service) UpdateShow(ctx context.Context, sh Show) error {
	if err := s.sr.UpdateShow(ctx, repository.UpdateShowParams{
		Title:         sh.Title,
		Cover:         sh.Cover,
		HasNewEpisode: sh.HasNewEpisode,
		Category: sql.NullString{
			String: sh.Category,
			Valid:  true,
		},
		Description: sql.NullString{
			String: sh.Description,
			Valid:  len(sh.Description) > 0,
		},
		ID: sh.ID,
	}); err != nil {
		return fmt.Errorf("could not update show with id=%s:%w", sh.ID, err)
	}
	return nil
}

// DeleteShowByID ..
func (s *Service) DeleteShowByID(ctx context.Context, id uuid.UUID) error {
	if err := s.sr.DeleteShowByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete show with id=%s:%w", id, err)
	}

	return nil
}

// AddEpisode ..
func (s *Service) AddEpisode(ctx context.Context, ep Episode) error {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return fmt.Errorf("could not add parse date from string: %w", err)
	}

	if err = s.sr.AddEpisode(ctx, repository.AddEpisodeParams{
		ShowID:        ep.ShowID,
		EpisodeNumber: ep.EpisodeNumber,
		Cover: sql.NullString{
			String: ep.Cover,
			Valid:  true,
		},
		Title: ep.Title,
		Description: sql.NullString{
			String: ep.Description,
			Valid:  true,
		},
		ReleaseDate: sql.NullTime{
			Time:  rDate,
			Valid: true,
		},
	}); err != nil {
		return fmt.Errorf("could not add episode with show_id=%s, episodeNumber=%v: %w", ep.ShowID, ep.EpisodeNumber, err)
	}

	return nil
}

// UpdateEpisode ..
func (s *Service) UpdateEpisode(ctx context.Context, ep Episode) error {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return fmt.Errorf("could not add parse date from string: %w", err)
	}

	if err = s.sr.UpdateEpisode(ctx, repository.UpdateEpisodeParams{
		ShowID:        ep.ShowID,
		EpisodeNumber: ep.EpisodeNumber,
		Cover: sql.NullString{
			String: ep.Cover,
			Valid:  true,
		},
		Title: ep.Title,
		Description: sql.NullString{
			String: ep.Title,
			Valid:  true,
		},
		ReleaseDate: sql.NullTime{
			Time:  rDate,
			Valid: true,
		},
		ID: ep.ID,
	}); err != nil {
		return fmt.Errorf("could not update episode with id=%s:%w", ep.ID, err)
	}
	return nil
}

// DeleteEpisodeByID ..
func (s *Service) DeleteEpisodeByID(ctx context.Context, showID, episodeID uuid.UUID) error {
	if err := s.sr.DeleteEpisodeByID(ctx, repository.DeleteEpisodeByIDParams{
		ID:     episodeID,
		ShowID: showID,
	}); err != nil {
		return fmt.Errorf("could not delete episode with id=%s:%w", episodeID, err)
	}

	return nil
}
