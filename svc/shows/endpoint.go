package shows

import (
	"context"
	"errors"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/utils"
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

		AddEpisode               endpoint.Endpoint
		DeleteEpisodeByID        endpoint.Endpoint
		GetActivatedUserEpisodes endpoint.Endpoint
		GetEpisodeByID           endpoint.Endpoint
		GetEpisodesByShowID      endpoint.Endpoint
		UpdateEpisode            endpoint.Endpoint

		RateEpisode            endpoint.Endpoint
		ReviewEpisode          endpoint.Endpoint
		GetReviewsList         endpoint.Endpoint
		GetReviewsListByUserID endpoint.Endpoint
		DeleteReviewByID       endpoint.Endpoint
		LikeDislikeEpisode     endpoint.Endpoint

		AddClapsForShow endpoint.Endpoint

		SendTipsToReviewAuthor endpoint.Endpoint
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
		GetActivatedUserEpisodes(ctx context.Context, userID uuid.UUID, page, itemsPerPage int32) ([]Episode, error)
		GetEpisodesByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
		GetEpisodeByID(ctx context.Context, showID, episodeID, userID uuid.UUID) (Episode, error)
		UpdateEpisode(ctx context.Context, ep Episode) error

		RateEpisode(ctx context.Context, episodeID, userID uuid.UUID, rating int32) error
		ReviewEpisode(ctx context.Context, episodeID, userID uuid.UUID, username string, rating int32, title, review string) error
		GetReviewsList(ctx context.Context, episodeID uuid.UUID, limit, offset int32, currentUserID uuid.UUID) ([]Review, error)
		GetReviewsListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32, currentUserID uuid.UUID) ([]Review, error)
		DeleteReviewByID(ctx context.Context, id uuid.UUID) error
		LikeDislikeEpisodeReview(ctx context.Context, id, uid uuid.UUID, ratingType ReviewRatingType) error

		AddClapsForShow(ctx context.Context, showID, userID uuid.UUID) error

		SendTipsToReviewAuthor(ctx context.Context, reviewID, uid uuid.UUID, amount float64) error
	}

	// GetShowChallengesRequest struct
	GetShowChallengesRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`
		utils.PaginationRequest
	}

	// GetShowsByCategoryRequest struct
	GetShowsByCategoryRequest struct {
		Category string `json:"category"`
		utils.PaginationRequest
	}

	// AddShowRequest struct
	AddShowRequest struct {
		Title          string `json:"title,omitempty" validate:"required,gt=0"`
		Cover          string `json:"cover,omitempty" validate:"required,gt=0"`
		HasNewEpisode  bool   `json:"has_new_episode,omitempty"`
		Category       string `json:"category,omitempty"`
		Description    string `json:"description,omitempty"`
		RealmsTitle    string `json:"realms_title,omitempty"`
		RealmsSubtitle string `json:"realms_subtitle,omitempty"`
		Watch          string `json:"watch,omitempty"`
	}

	// UpdateShowRequest struct
	UpdateShowRequest struct {
		ID             string `json:"id,omitempty" validate:"required,uuid"`
		Title          string `json:"title,omitempty" validate:"required"`
		Cover          string `json:"cover,omitempty" validate:"required"`
		HasNewEpisode  bool   `json:"has_new_episode,omitempty"`
		Category       string `json:"category,omitempty"`
		Description    string `json:"description,omitempty"`
		RealmsTitle    string `json:"realms_title,omitempty"`
		RealmsSubtitle string `json:"realms_subtitle,omitempty"`
		Watch          string `json:"watch,omitempty"`
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
		HintText                string `json:"hint_text,omitempty"`
		Watch                   string `json:"watch,omitempty"`
	}

	// UpdateEpisodeRequest struct
	UpdateEpisodeRequest struct {
		ID                      string `json:"id" validate:"required,uuid"`
		ShowID                  string `json:"show_id" validate:"required,uuid"`
		SeasonID                string `json:"season_id" validate:"required,uuid"`
		EpisodeNumber           int32  `json:"episode_number"`
		Cover                   string `json:"cover,omitempty"`
		Title                   string `json:"title" validate:"required,gt=0"`
		Description             string `json:"description,omitempty"`
		ReleaseDate             string `json:"release_date" validate:"datetime=2006-01-02T15:04:05Z"`
		ChallengeID             string `json:"challenge_id,omitempty"`
		VerificationChallengeID string `json:"verification_challenge_id,omitempty"`
		HintText                string `json:"hint_text,omitempty"`
		Watch                   string `json:"watch,omitempty"`
	}

	// GetEpisodesByShowIDRequest struct
	GetEpisodesByShowIDRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`
		utils.PaginationRequest
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
		Rating    int32  `json:"rating" validate:"required,gte=1,lte=10"`
	}

	// ReviewEpisodeRequest struct
	ReviewEpisodeRequest struct {
		EpisodeID string `json:"episode_id,omitempty" validate:"required,uuid"`
		Rating    int32  `json:"rating" validate:"required,gte=1,lte=10"`
		Title     string `json:"title" validate:"required"`
		Review    string `json:"review" validate:"required"`
	}

	// SendTipsRequest struct
	SendTipsRequest struct {
		ReviewID string  `json:"review_id" validate:"required,uuid"`
		Amount   float64 `json:"amount" validate:"required"`
	}

	// GetReviewsListRequest struct
	GetReviewsListRequest struct {
		EpisodeID string `json:"episode_id" validate:"required,uuid"`
		utils.PaginationRequest
	}

	// LikeDislikeEpisodeRequest struct
	LikeDislikeEpisodeRequest struct {
		ReviewID string `json:"review_id" validate:"required,uuid"`
		Param    string `json:"rating_type" validate:"required,oneof=like dislike"`
	}
)

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

		AddEpisode:               MakeAddEpisodeEndpoint(s, validateFunc),
		DeleteEpisodeByID:        MakeDeleteEpisodeByIDEndpoint(s, validateFunc),
		GetActivatedUserEpisodes: MakeGetActivatedUserEpisodesEndpoint(s, validateFunc),
		GetEpisodeByID:           MakeGetEpisodeByIDEndpoint(s, validateFunc),
		GetEpisodesByShowID:      MakeGetEpisodesByShowIDEndpoint(s, validateFunc),
		UpdateEpisode:            MakeUpdateEpisodeEndpoint(s, validateFunc),

		RateEpisode:            MakeRateEpisodeEndpoint(s, validateFunc),
		ReviewEpisode:          MakeReviewEpisodeEndpoint(s, validateFunc),
		GetReviewsList:         MakeGetReviewsListEndpoint(s, validateFunc),
		GetReviewsListByUserID: MakeGetReviewsListByUserIDEndpoint(s, validateFunc),
		DeleteReviewByID:       MakeDeleteReviewByIDEndpoint(s),
		LikeDislikeEpisode:     MakeLikeDislikeEpisodeEndpoint(s, validateFunc),

		AddClapsForShow: MakeAddClapsForShowEndpoint(s),

		SendTipsToReviewAuthor: MakeSendTipsToReviewAuthorEndpoint(s, validateFunc),
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
			e.GetActivatedUserEpisodes = mdw(e.GetActivatedUserEpisodes)
			e.GetEpisodeByID = mdw(e.GetEpisodeByID)
			e.GetEpisodesByShowID = mdw(e.GetEpisodesByShowID)
			e.UpdateEpisode = mdw(e.UpdateEpisode)

			e.RateEpisode = mdw(e.RateEpisode)
			e.ReviewEpisode = mdw(e.ReviewEpisode)
			e.GetReviewsList = mdw(e.GetReviewsList)
			e.GetReviewsListByUserID = mdw(e.GetReviewsListByUserID)
			e.DeleteReviewByID = mdw(e.DeleteReviewByID)
			e.LikeDislikeEpisode = mdw(e.LikeDislikeEpisode)

			e.AddClapsForShow = mdw(e.AddClapsForShow)

			e.SendTipsToReviewAuthor = mdw(e.SendTipsToReviewAuthor)
		}
	}

	return e
}

// MakeGetShowsEndpoint ...
func MakeGetShowsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(utils.PaginationRequest)
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(AddShowRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddShow(ctx, Show{
			Title:          req.Title,
			Cover:          req.Cover,
			HasNewEpisode:  req.HasNewEpisode,
			Category:       req.Category,
			Description:    req.Description,
			RealmsTitle:    req.RealmsTitle,
			RealmsSubtitle: req.RealmsSubtitle,
			Watch:          req.Watch,
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(UpdateShowRequest)

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.UpdateShow(ctx, Show{
			ID:             id,
			Title:          req.Title,
			Cover:          req.Cover,
			HasNewEpisode:  req.HasNewEpisode,
			Category:       req.Category,
			Description:    req.Description,
			RealmsTitle:    req.RealmsTitle,
			RealmsSubtitle: req.RealmsSubtitle,
			Watch:          req.Watch,
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

		payload := Episode{
			ShowID:        showID,
			SeasonID:      seasonID,
			EpisodeNumber: req.EpisodeNumber,
			Cover:         req.Cover,
			Title:         req.Title,
			Description:   req.Description,
			ReleaseDate:   req.ReleaseDate,
			HintText:      req.HintText,
			Watch:         req.Watch,
		}

		if req.ChallengeID != "" {
			challengeID, err := uuid.Parse(req.ChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get challenge id: %w", err)
			}
			if challengeID != uuid.Nil {
				payload.ChallengeID = &challengeID
			}
		}

		if req.VerificationChallengeID != "" {
			verificationChallengeID, err := uuid.Parse(req.VerificationChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get verification challenge id: %w", err)
			}
			if verificationChallengeID != uuid.Nil {
				payload.VerificationChallengeID = &verificationChallengeID
			}
		}

		resp, err := s.AddEpisode(ctx, payload)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeUpdateEpisodeEndpoint ...
func MakeUpdateEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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

		payload := Episode{
			ID:            id,
			ShowID:        showID,
			SeasonID:      seasonID,
			EpisodeNumber: req.EpisodeNumber,
			Cover:         req.Cover,
			Title:         req.Title,
			Description:   req.Description,
			ReleaseDate:   req.ReleaseDate,
			HintText:      req.HintText,
			Watch:         req.Watch,
		}

		if req.ChallengeID != "" {
			challengeID, err := uuid.Parse(req.ChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get challenge id: %w", err)
			}
			if challengeID != uuid.Nil {
				payload.ChallengeID = &challengeID
			}
		}

		if req.VerificationChallengeID != "" {
			verificationChallengeID, err := uuid.Parse(req.VerificationChallengeID)
			if err != nil {
				return nil, fmt.Errorf("could not get verification challenge id: %w", err)
			}
			if verificationChallengeID != uuid.Nil {
				payload.VerificationChallengeID = &verificationChallengeID
			}
		}

		if err := s.UpdateEpisode(ctx, payload); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteEpisodeByIDEndpoint ...
func MakeDeleteEpisodeByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

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

		resp, err := s.GetEpisodeByID(ctx, showID, episodeID, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetActivatedUserEpisodesEndpoint ...
func MakeGetActivatedUserEpisodesEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(utils.PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetActivatedUserEpisodes(ctx, uid, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetEpisodesByShowIDEndpoint ...
func MakeGetEpisodesByShowIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(GetEpisodesByShowIDRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		resp, err := s.GetEpisodesByShowID(ctx, showID, uid, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddSeasonEndpoint ...
func MakeAddSeasonEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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

// MakeReviewEpisodeEndpoint ...
func MakeReviewEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		req := request.(ReviewEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		if err := s.ReviewEpisode(ctx, episodeID, uid, username, req.Rating, req.Title, req.Review); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetReviewsListEndpoint ...
func MakeGetReviewsListEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(GetReviewsListRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetReviewsList(ctx, episodeID, req.Limit(), req.Offset(), uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetReviewsListByUserIDEndpoint ...
func MakeGetReviewsListByUserIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(utils.PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetReviewsListByUserID(ctx, uid, req.Limit(), req.Offset(), uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteReviewByIDEndpoint ...
func MakeDeleteReviewByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get review id: %w", err)
		}

		err = s.DeleteReviewByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeLikeDislikeEpisodeEndpoint ...
func MakeLikeDislikeEpisodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(LikeDislikeEpisodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		var ratingType ReviewRatingType
		switch req.Param {
		case "like":
			ratingType = LikeReview
		case "dislike":
			ratingType = DislikeReview
		default:
			return nil, fmt.Errorf("undefined rating type: %s", req.Param)
		}

		reviewID, err := uuid.Parse(req.ReviewID)
		if err != nil {
			return nil, fmt.Errorf("%w review id: %v", ErrInvalidParameter, err)
		}

		if err := s.LikeDislikeEpisodeReview(ctx, reviewID, uid, ratingType); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeAddClapsForShowEndpoint ...
func MakeAddClapsForShowEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		showID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w episode id: %v", ErrInvalidParameter, err)
		}

		if err := s.AddClapsForShow(ctx, showID, uid); err != nil {
			if errors.Is(err, ErrMaxClaps) {
				return false, nil
			}
			return nil, err
		}

		return true, nil
	}
}

// MakeSendTipsToReviewAuthorEndpoint ...
func MakeSendTipsToReviewAuthorEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(SendTipsRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		reviewID, err := uuid.Parse(req.ReviewID)
		if err != nil {
			return nil, fmt.Errorf("%w review id: %v", ErrInvalidParameter, err)
		}

		if err := s.SendTipsToReviewAuthor(ctx, reviewID, uid, req.Amount); err != nil {
			if errors.Is(err, ErrMaxClaps) {
				return false, nil
			}
			return nil, err
		}

		return true, nil
	}
}
