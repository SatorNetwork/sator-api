package shows

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/internal/db"
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
		SeasonNumber int32     `json:"season_number"`
		Episodes     []Episode `json:"episodes"`
		ShowID       uuid.UUID `json:"show_id"`
	}

	Episode struct {
		ID                      uuid.UUID `json:"id"`
		ShowID                  uuid.UUID `json:"show_id"`
		SeasonID                uuid.UUID `json:"season_id"`
		SeasonNumber            int32     `json:"season_number"`
		EpisodeNumber           int32     `json:"episode_number"`
		Cover                   string    `json:"cover"`
		Title                   string    `json:"title"`
		Description             string    `json:"description"`
		ReleaseDate             string    `json:"release_date"`
		ChallengeID             uuid.UUID `json:"challenge_id"`
		VerificationChallengeID uuid.UUID `json:"verification_challenge_id"`
		Rating                  float64   `json:"rating"`
		RatingsCount            int64     `json:"ratings_count"`
	}

	showsRepository interface {
		// Shows
		AddShow(ctx context.Context, arg repository.AddShowParams) (repository.Show, error)
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetShows(ctx context.Context, arg repository.GetShowsParams) ([]repository.Show, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.Show, error)
		GetShowsByCategory(ctx context.Context, arg repository.GetShowsByCategoryParams) ([]repository.Show, error)
		UpdateShow(ctx context.Context, arg repository.UpdateShowParams) error

		// Seasons
		AddSeason(ctx context.Context, arg repository.AddSeasonParams) (repository.Season, error)
		DeleteSeasonByID(ctx context.Context, arg repository.DeleteSeasonByIDParams) error
		GetSeasonByID(ctx context.Context, id uuid.UUID) (repository.Season, error)
		GetSeasonsByShowID(ctx context.Context, arg repository.GetSeasonsByShowIDParams) ([]repository.Season, error)

		// Episodes
		AddEpisode(ctx context.Context, arg repository.AddEpisodeParams) (repository.Episode, error)
		GetEpisodeByID(ctx context.Context, id uuid.UUID) (repository.GetEpisodeByIDRow, error)
		GetEpisodesByShowID(ctx context.Context, arg repository.GetEpisodesByShowIDParams) ([]repository.GetEpisodesByShowIDRow, error)
		DeleteEpisodeByID(ctx context.Context, id uuid.UUID) error
		UpdateEpisode(ctx context.Context, arg repository.UpdateEpisodeParams) error

		// Episodes rating
		GetEpisodeRatingByID(ctx context.Context, episodeID uuid.UUID) (repository.GetEpisodeRatingByIDRow, error)
		RateEpisode(ctx context.Context, arg repository.RateEpisodeParams) error
	}

	// Challenges service client
	challengesClient interface {
		GetListByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
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
func (s *Service) GetShowChallenges(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error) {
	challenges, err := s.chc.GetListByShowID(ctx, showID, userID, limit, offset)
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
		if _, ok := episodesPerSeasons[e.SeasonID.UUID.String()]; ok {
			episodesPerSeasons[e.SeasonID.UUID.String()] = append(episodesPerSeasons[e.SeasonID.UUID.String()], castRowsToEpisode(e))
		} else {
			episodesPerSeasons[e.SeasonID.UUID.String()] = []Episode{castRowsToEpisode(e)}
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
			SeasonNumber: s.SeasonNumber,
			Title:        fmt.Sprintf("Season %d", s.SeasonNumber),
			Episodes:     episodes[s.ID.String()],
			ShowID:       s.ShowID,
		})
	}
	return result
}

// GetEpisodeByID returns episode with provided id.
func (s *Service) GetEpisodeByID(ctx context.Context, showID, episodeID uuid.UUID) (interface{}, error) {
	episode, err := s.sr.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get episode with id=%s: %w", episodeID, err)
	}

	avgRating, ratingsCount, err := s.getAverageEpisodesRatingByID(ctx, episodeID)
	if err != nil {
		return nil, fmt.Errorf("could not get avarage episoderating with id=%s: %w", episodeID, err)
	}

	return castRowToEpisode(episode, avgRating, ratingsCount), nil
}

// Cast repository.GetEpisodeByIDRow to service Episode structure
func castRowToEpisode(source repository.GetEpisodeByIDRow, rating float64, ratingsCount int64) Episode {
	return Episode{
		ID:                      source.ID,
		ShowID:                  source.ShowID,
		EpisodeNumber:           source.EpisodeNumber,
		SeasonID:                source.SeasonID.UUID,
		SeasonNumber:            source.SeasonNumber,
		Cover:                   source.Cover.String,
		Title:                   source.Title,
		Description:             source.Description.String,
		ReleaseDate:             source.ReleaseDate.Time.String(),
		ChallengeID:             source.ChallengeID.UUID,
		VerificationChallengeID: source.VerificationChallengeID.UUID,
		Rating:                  rating,
		RatingsCount:            ratingsCount,
	}
}

// Cast repository.GetEpisodesByShowIDRow to service Episode structure
func castRowsToEpisode(source repository.GetEpisodesByShowIDRow) Episode {
	return Episode{
		ID:                      source.ID,
		ShowID:                  source.ShowID,
		EpisodeNumber:           source.EpisodeNumber,
		SeasonID:                source.SeasonID.UUID,
		SeasonNumber:            source.SeasonNumber,
		Cover:                   source.Cover.String,
		Title:                   source.Title,
		Description:             source.Description.String,
		ReleaseDate:             source.ReleaseDate.Time.String(),
		ChallengeID:             source.ChallengeID.UUID,
		VerificationChallengeID: source.VerificationChallengeID.UUID,
		Rating:                  source.AvgRating,
		RatingsCount:            source.Ratings,
	}
}

// Cast repository.Episode to service Episode structure
func castToEpisode(source repository.Episode, seasonNumber int32) Episode {
	return Episode{
		ID:                      source.ID,
		ShowID:                  source.ShowID,
		EpisodeNumber:           source.EpisodeNumber,
		SeasonID:                source.SeasonID.UUID,
		SeasonNumber:            seasonNumber,
		Cover:                   source.Cover.String,
		Title:                   source.Title,
		Description:             source.Description.String,
		ReleaseDate:             source.ReleaseDate.Time.String(),
		ChallengeID:             source.ChallengeID.UUID,
		VerificationChallengeID: source.VerificationChallengeID.UUID,
	}
}

// AddShow ...
func (s *Service) AddShow(ctx context.Context, sh Show) (Show, error) {
	show, err := s.sr.AddShow(ctx, repository.AddShowParams{
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
	})
	if err != nil {
		return Show{}, fmt.Errorf("could not add show with title=%s: %w", sh.Title, err)
	}

	return castToShow(show), nil
}

// UpdateShow ...
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
func (s *Service) AddEpisode(ctx context.Context, ep Episode) (Episode, error) {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return Episode{}, fmt.Errorf("could not add parse date from string: %w", err)
	}

	episode, err := s.sr.AddEpisode(ctx, repository.AddEpisodeParams{
		ShowID:        ep.ShowID,
		SeasonID:      uuid.NullUUID{UUID: ep.SeasonID, Valid: ep.SeasonID != uuid.Nil},
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
		ChallengeID:             uuid.NullUUID{UUID: ep.ChallengeID, Valid: ep.ChallengeID != uuid.Nil},
		VerificationChallengeID: uuid.NullUUID{UUID: ep.VerificationChallengeID, Valid: ep.VerificationChallengeID != uuid.Nil},
	})
	if err != nil {
		return Episode{}, fmt.Errorf("could not add episode #%d for show_id=%s: %w", ep.EpisodeNumber, ep.ShowID.String(), err)
	}

	episodeByID, err := s.sr.GetEpisodeByID(ctx, episode.ID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get episode with id=%s: %w", episode.ID, err)
	}

	return castToEpisode(episode, episodeByID.SeasonNumber), nil
}

// UpdateEpisode ..
func (s *Service) UpdateEpisode(ctx context.Context, ep Episode) error {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return fmt.Errorf("could not add parse date from string: %w", err)
	}

	if err = s.sr.UpdateEpisode(ctx, repository.UpdateEpisodeParams{
		ID:            ep.ID,
		ShowID:        ep.ShowID,
		SeasonID:      uuid.NullUUID{UUID: ep.SeasonID, Valid: ep.SeasonID != uuid.Nil},
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
		ChallengeID:             uuid.NullUUID{UUID: ep.ChallengeID, Valid: ep.ChallengeID != uuid.Nil},
		VerificationChallengeID: uuid.NullUUID{UUID: ep.VerificationChallengeID, Valid: ep.VerificationChallengeID != uuid.Nil},
	}); err != nil {
		return fmt.Errorf("could not update episode with id=%s:%w", ep.ID, err)
	}

	return nil
}

// DeleteEpisodeByID ...
func (s *Service) DeleteEpisodeByID(ctx context.Context, showID, episodeID uuid.UUID) error {
	if err := s.sr.DeleteEpisodeByID(ctx, episodeID); err != nil {
		return fmt.Errorf("could not delete episode with id=%s:%w", episodeID, err)
	}

	return nil
}

// AddSeason ...
func (s *Service) AddSeason(ctx context.Context, ss Season) (Season, error) {
	season, err := s.sr.AddSeason(ctx, repository.AddSeasonParams{
		ShowID:       ss.ShowID,
		SeasonNumber: ss.SeasonNumber,
	})
	if err != nil {
		return Season{}, fmt.Errorf("could not add season with title=%s: %w", ss.Title, err)
	}

	return castToSeason(season), nil
}

// Cast repository.Season to service Season structure
func castToSeason(source repository.Season) Season {
	return Season{
		ID:           source.ID,
		Title:        fmt.Sprintf("Season %d", source.SeasonNumber),
		SeasonNumber: source.SeasonNumber,
		ShowID:       source.ShowID,
	}
}

// DeleteSeasonByID ...
func (s *Service) DeleteSeasonByID(ctx context.Context, showID, seasonID uuid.UUID) error {
	if err := s.sr.DeleteSeasonByID(ctx, repository.DeleteSeasonByIDParams{
		ID:     seasonID,
		ShowID: showID,
	}); err != nil {
		return fmt.Errorf("could not delete season with id=%s:%w", seasonID, err)
	}

	return nil
}

// getAverageEpisodesRatingByID returns average episode rating.
func (s *Service) getAverageEpisodesRatingByID(ctx context.Context, episodeID uuid.UUID) (float64, int64, error) {
	rating, err := s.sr.GetEpisodeRatingByID(ctx, episodeID)
	if err != nil && !db.IsNotFoundError(err) {
		return 0, 0, fmt.Errorf("could not get average episode rating by ID= %v: %w", episodeID, err)
	}

	return rating.AvgRating, rating.Ratings, nil
}

// RateEpisode ...
func (s *Service) RateEpisode(ctx context.Context, episodeID, userID uuid.UUID, rating int32) error {
	err := s.sr.RateEpisode(ctx, repository.RateEpisodeParams{
		EpisodeID: episodeID,
		UserID:    userID,
		Rating:    rating,
	})
	if err != nil {
		if db.IsDuplicateError(err) {
			return fmt.Errorf("you already rate this episode")
		}
		return fmt.Errorf("could not rate episode with episodeID=%s: %w", episodeID, err)
	}

	return nil
}
