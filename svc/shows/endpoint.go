package shows

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddShow            endpoint.Endpoint
		DeleteShowByID     endpoint.Endpoint
		GetShows           endpoint.Endpoint
		GetShowChallenges  endpoint.Endpoint
		GetShowByID        endpoint.Endpoint
		GetShowsByCategory endpoint.Endpoint
		UpdateShow         endpoint.Endpoint

		AddSeason        endpoint.Endpoint
		DeleteSeasonByID endpoint.Endpoint

		AddEpisode          endpoint.Endpoint
		DeleteEpisodeByID   endpoint.Endpoint
		GetEpisodeByID      endpoint.Endpoint
		GetEpisodesByShowID endpoint.Endpoint
		UpdateEpisode       endpoint.Endpoint

		RateEpisode endpoint.Endpoint
	}

	service interface {
		AddShow(ctx context.Context, sh Show) (Show, error)
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetShows(ctx context.Context, page, itemsPerPage int32) (interface{}, error)
		GetShowChallenges(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetShowsByCategory(ctx context.Context, category string, limit, offset int32) (interface{}, error)
		UpdateShow(ctx context.Context, sh Show) error

		AddSeason(ctx context.Context, ss Season) (Season, error)
		DeleteSeasonByID(ctx context.Context, showID, seasonID uuid.UUID) error

		AddEpisode(ctx context.Context, ep Episode) (Episode, error)
		DeleteEpisodeByID(ctx context.Context, showId, episodeId uuid.UUID) error
		GetEpisodesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
		GetEpisodeByID(ctx context.Context, showID, episodeID uuid.UUID) (interface{}, error)
		UpdateEpisode(ctx context.Context, ep Episode) error

		RateEpisode(ctx context.Context, episodeID, userID uuid.UUID, rating int32) error
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}

	// GetShowChallengesRequest struct
	GetShowChallengesRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`
		PaginationRequest
	}

	// GetShowsByCategoryRequest struct
	GetShowsByCategoryRequest struct {
		Category string `json:"category"`
		PaginationRequest
	}

	// AddShowRequest struct
	AddShowRequest struct {
		Title         string `json:"title,omitempty" validate:"required,gt=0"`
		Cover         string `json:"cover,omitempty" validate:"required,gt=0"`
		HasNewEpisode bool   `json:"has_new_episode,omitempty"`
		Category      string `json:"category,omitempty"`
		Description   string `json:"description,omitempty"`
	}

	// UpdateShowRequest struct
	UpdateShowRequest struct {
		ID            string `json:"id,omitempty" validate:"required,uuid"`
		Title         string `json:"title,omitempty" validate:"required"`
		Cover         string `json:"cover,omitempty" validate:"required"`
		HasNewEpisode bool   `json:"has_new_episode,omitempty"`
		Category      string `json:"category,omitempty"`
		Description   string `json:"description,omitempty"`
	}

	// GetEpisodeByIDRequest struct
	GetEpisodeByIDRequest struct {
		ShowID    string `json:"show_id" validate:"required,uuid"`
		EpisodeID string `json:"episode_id" validate:"required,uuid"`
	}

	// DeleteEpisodeByIDRequest struct
	DeleteEpisodeByIDRequest struct {
		ShowID    string `json:"show_id" validate:"required,uuid"`
		EpisodeID string `json:"episode_id" validate:"required,uuid"`
	}

	// AddEpisodeRequest struct
	AddEpisodeRequest struct {
		ShowID                  string `json:"show_id" validate:"required,uuid"`
		SeasonID                string `json:"season_id" validate:"required,uuid"`
		EpisodeNumber           int32  `json:"episode_number"`
		Cover                   string `json:"cover,omitempty"`
		Title                   string `json:"title" validate:"required,gt=0"`
		Description             string `json:"description,omitempty"`
		ReleaseDate             string `json:"release_date,omitempty"`
		ChallengeID             string `json:"challenge_id,omitempty"`
		VerificationChallengeID string `json:"verification_challenge_id,omitempty"`
	}

	// UpdateEpisodeRequest struct
	UpdateEpisodeRequest struct {
		ID                      string `json:"id" validate:"required,uuid"`
		ShowID                  string `json:"show_id" validate:"required,uuid"`
		SeasonID                string `json:"season_id" validate:"required,uuid"`
		EpisodeNumber           int32  `json:"episode_number"`
		Cover                   string `json:"cover,omitempty"`
		Title                   string `json:"title" validate:"required,gt=0"`
		Description             string `json:"description" validate:"required,gt=0"`
		ReleaseDate             string `json:"release_date" validate:"datetime=2006-01-02T15:04:05Z"`
		ChallengeID             string `json:"challenge_id,omitempty"`
		VerificationChallengeID string `json:"verification_challenge_id,omitempty"`
	}

	// GetEpisodesByShowIDRequest struct
	GetEpisodesByShowIDRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`
		PaginationRequest
	}

	// AddSeasonRequest struct
	AddSeasonRequest struct {
		ShowID       string `json:"show_id" validate:"required,uuid"`
		SeasonNumber int32  `json:"season_number"`
	}

	// DeleteSeasonByIDRequest struct
	DeleteSeasonByIDRequest struct {
		SeasonID string `json:"season_id" validate:"required,uuid"`
		ShowID   string `json:"show_id" validate:"required,uuid"`
	}

	// RateEpisodeRequest struct
	RateEpisodeRequest struct {
		EpisodeID string `json:"episode_id" validate:"required,uuid"`
		Rating    int32  `json:"rating" validate:"required"`
	}
)

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}
	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}
	return 0
}

// MakeEndpoints ...
func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddShow:            MakeAddShowEndpoint(s, validateFunc),
		DeleteShowByID:     MakeDeleteShowByIDEndpoint(s),
		GetShows:           MakeGetShowsEndpoint(s, validateFunc),
		GetShowChallenges:  MakeGetShowChallengesEndpoint(s, validateFunc),
		GetShowByID:        MakeGetShowByIDEndpoint(s),
		GetShowsByCategory: MakeGetShowsByCategoryEndpoint(s, validateFunc),
		UpdateShow:         MakeUpdateShowEndpoint(s),

		AddSeason:        MakeAddSeasonEndpoint(s, validateFunc),
		DeleteSeasonByID: MakeDeleteSeasonByIDEndpoint(s, validateFunc),

		AddEpisode:          MakeAddEpisodeEndpoint(s, validateFunc),
		DeleteEpisodeByID:   MakeDeleteEpisodeByIDEndpoint(s, validateFunc),
		GetEpisodeByID:      MakeGetEpisodeByIDEndpoint(s, validateFunc),
		GetEpisodesByShowID: MakeGetEpisodesByShowIDEndpoint(s, validateFunc),
		UpdateEpisode:       MakeUpdateEpisodeEndpoint(s, validateFunc),

		RateEpisode: MakeRateEpisodeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddShow = mdw(e.AddShow)
			e.DeleteShowByID = mdw(e.DeleteShowByID)
			e.GetShows = mdw(e.GetShows)
			e.GetShowChallenges = mdw(e.GetShowChallenges)
			e.GetShowByID = mdw(e.GetShowByID)
			e.GetShowsByCategory = mdw(e.GetShowsByCategory)
			e.UpdateShow = mdw(e.UpdateShow)

			e.AddSeason = mdw(e.AddSeason)
			e.DeleteSeasonByID = mdw(e.DeleteSeasonByID)

			e.AddEpisode = mdw(e.AddEpisode)
			e.DeleteEpisodeByID = mdw(e.DeleteEpisodeByID)
			e.GetEpisodeByID = mdw(e.GetEpisodeByID)
			e.GetEpisodesByShowID = mdw(e.GetEpisodesByShowID)
			e.UpdateEpisode = mdw(e.UpdateEpisode)

			e.RateEpisode = mdw(e.RateEpisode)
		}
	}

	return e
}

// MakeGetShowsEndpoint ...
func MakeGetShowsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetShows(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetShowChallengesEndpoint ...
func MakeGetShowChallengesEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(GetShowChallengesRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetShowChallenges(ctx, showID, uid, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetShowByIDEndpoint ...
func MakeGetShowByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetShowByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetShowsByCategoryEndpoint ...
func MakeGetShowsByCategoryEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetShowsByCategoryRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		if req.Category != "" {
			resp, err := s.GetShowsByCategory(ctx, req.Category, req.Limit(), req.Offset())
			if err != nil {
				return nil, err
			}

			return resp, nil
		}

		resp, err := s.GetShows(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddShowEndpoint ...
func MakeAddShowEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddShowRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddShow(ctx, Show{
			Title:         req.Title,
			Cover:         req.Cover,
			HasNewEpisode: req.HasNewEpisode,
			Category:      req.Category,
			Description:   req.Description,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeUpdateShowEndpoint ...
func MakeUpdateShowEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateShowRequest)

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.UpdateShow(ctx, Show{
			ID:            id,
			Title:         req.Title,
			Cover:         req.Cover,
			HasNewEpisode: req.HasNewEpisode,
			Category:      req.Category,
			Description:   req.Description,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteShowByIDEndpoint ...
func MakeDeleteShowByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.DeleteShowByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeAddEpisodeEndpoint ...
func MakeAddEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		seasonID, err := uuid.Parse(req.SeasonID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		var challengeID, verificationChallengeID uuid.UUID
		if req.ChallengeID != "" {
			challengeID, err = uuid.Parse(req.ChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get challenge id: %w", err)
			}
		}

		if req.VerificationChallengeID != "" {
			verificationChallengeID, err = uuid.Parse(req.VerificationChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get verification challenge id: %w", err)
			}
		}

		resp, err := s.AddEpisode(ctx, Episode{
			ShowID:                  showID,
			SeasonID:                seasonID,
			ChallengeID:             challengeID,
			VerificationChallengeID: verificationChallengeID,
			EpisodeNumber:           req.EpisodeNumber,
			Cover:                   req.Cover,
			Title:                   req.Title,
			Description:             req.Description,
			ReleaseDate:             req.ReleaseDate,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeUpdateEpisodeEndpoint ...
func MakeUpdateEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		seasonID, err := uuid.Parse(req.SeasonID)
		if err != nil {
			return nil, fmt.Errorf("could not get season id: %w", err)
		}

		var challengeID, verificationChallengeID uuid.UUID
		if req.ChallengeID != "" {
			challengeID, err = uuid.Parse(req.ChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get challenge id: %w", err)
			}
		}

		if req.VerificationChallengeID != "" {
			verificationChallengeID, err = uuid.Parse(req.VerificationChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get verification challenge id: %w", err)
			}
		}

		err = s.UpdateEpisode(ctx, Episode{
			ID:                      id,
			ShowID:                  showID,
			SeasonID:                seasonID,
			EpisodeNumber:           req.EpisodeNumber,
			Cover:                   req.Cover,
			Title:                   req.Title,
			Description:             req.Description,
			ReleaseDate:             req.ReleaseDate,
			ChallengeID:             challengeID,
			VerificationChallengeID: verificationChallengeID,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteEpisodeByIDEndpoint ...
func MakeDeleteEpisodeByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteEpisodeByIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		if err := s.DeleteEpisodeByID(ctx, showID, episodeID); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetEpisodeByIDEndpoint ...
func MakeGetEpisodeByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetEpisodeByIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetEpisodeByID(ctx, showID, episodeID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetEpisodesByShowIDEndpoint ...
func MakeGetEpisodesByShowIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetEpisodesByShowIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		resp, err := s.GetEpisodesByShowID(ctx, showID, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddSeasonEndpoint ...
func MakeAddSeasonEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddSeasonRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		resp, err := s.AddSeason(ctx, Season{
			ShowID:       showID,
			SeasonNumber: req.SeasonNumber,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteSeasonByIDEndpoint ...
func MakeDeleteSeasonByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteSeasonByIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("%w show id: %v", ErrInvalidParameter, err)
		}

		seasonID, err := uuid.Parse(req.SeasonID)
		if err != nil {
			return nil, fmt.Errorf("%w season id: %v", ErrInvalidParameter, err)
		}

		err = s.DeleteSeasonByID(ctx, showID, seasonID)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeRateEpisodeEndpoint ...
func MakeRateEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(RateEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		err = s.RateEpisode(ctx, episodeID, uid, req.Rating)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
