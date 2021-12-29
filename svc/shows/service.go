package shows

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/SatorNetwork/sator-api/svc/profile"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/svc/shows/repository"

	"github.com/google/uuid"
)

const (
	defaultHintText = "Start watching to earn SAO"

	LikeReview    ReviewRatingType = 1
	DislikeReview ReviewRatingType = 2
)

type (
	// ReviewRatingType ...
	ReviewRatingType int32

	// Service struct
	Service struct {
		sr  showsRepository
		chc challengesClient
		pc  profileClient
		ac  authClient
	}

	// Show struct
	// Fields were rearranged to optimize memory usage.
	Show struct {
		ID             uuid.UUID `json:"id"`
		Title          string    `json:"title"`
		Cover          string    `json:"cover"`
		HasNewEpisode  bool      `json:"has_new_episode"`
		Category       string    `json:"category"`
		Description    string    `json:"description"`
		Claps          int64     `json:"claps"`
		RealmsTitle    string    `json:"realms_title"`
		RealmsSubtitle string    `json:"realms_subtitle"`
		Watch          string    `json:"watch"`
	}

	Season struct {
		ID           uuid.UUID `json:"id"`
		Title        string    `json:"title"`
		SeasonNumber int32     `json:"season_number"`
		Episodes     []Episode `json:"episodes"`
		ShowID       uuid.UUID `json:"show_id"`
	}

	Episode struct {
		ID                      uuid.UUID  `json:"id"`
		ShowID                  uuid.UUID  `json:"show_id"`
		SeasonID                uuid.UUID  `json:"season_id"`
		SeasonNumber            int32      `json:"season_number"`
		EpisodeNumber           int32      `json:"episode_number"`
		Cover                   string     `json:"cover"`
		Title                   string     `json:"title"`
		ShowTitle               string     `json:"show_title,omitempty"`
		Description             string     `json:"description"`
		ReleaseDate             string     `json:"release_date"`
		ChallengeID             *uuid.UUID `json:"challenge_id"`
		VerificationChallengeID *uuid.UUID `json:"verification_challenge_id"`
		Rating                  float64    `json:"rating"`
		RatingsCount            int64      `json:"ratings_count"`
		ActiveUsers             int32      `json:"active_users"`
		UserRewardsAmount       float64    `json:"user_rewards_amount"`
		TotalRewardsAmount      float64    `json:"total_rewards_amount"`
		HintText                string     `json:"hint_text"`
		Watch                   string     `json:"watch"`
	}

	// Review ...
	Review struct {
		ID         string `json:"id"`
		UserID     string `json:"user_id"`
		Username   string `json:"username"`
		UserAvatar string `json:"user_avatar"`
		Rating     int    `json:"rating"`
		Title      string `json:"title"`
		Review     string `json:"review"`
		Likes      int64  `json:"likes"`
		Dislikes   int64  `json:"dislikes"`
		IsLiked    bool   `json:"is_liked"`
		IsDisliked bool   `json:"is_disliked"`
		CreatedAt  string `json:"created_at"`
	}

	showsRepository interface {
		// Shows
		AddShow(ctx context.Context, arg repository.AddShowParams) (repository.Show, error)
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetShows(ctx context.Context, arg repository.GetShowsParams) ([]repository.Show, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.GetShowByIDRow, error)
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
		GetListEpisodesByIDs(ctx context.Context, episodeIds []uuid.UUID) ([]repository.GetListEpisodesByIDsRow, error)
		GetEpisodesByShowID(ctx context.Context, arg repository.GetEpisodesByShowIDParams) ([]repository.GetEpisodesByShowIDRow, error)
		DeleteEpisodeByID(ctx context.Context, id uuid.UUID) error
		UpdateEpisode(ctx context.Context, arg repository.UpdateEpisodeParams) error

		// Episodes rating
		GetEpisodeRatingByID(ctx context.Context, episodeID uuid.UUID) (repository.GetEpisodeRatingByIDRow, error)
		RateEpisode(ctx context.Context, arg repository.RateEpisodeParams) error
		DidUserRateEpisode(ctx context.Context, arg repository.DidUserRateEpisodeParams) (bool, error)

		// Episode reviews
		DidUserReviewEpisode(ctx context.Context, arg repository.DidUserReviewEpisodeParams) (bool, error)
		ReviewEpisode(ctx context.Context, arg repository.ReviewEpisodeParams) error
		ReviewsList(ctx context.Context, arg repository.ReviewsListParams) ([]repository.Rating, error)
		ReviewsListByUserID(ctx context.Context, arg repository.ReviewsListByUserIDParams) ([]repository.Rating, error)
		DeleteReview(ctx context.Context, id uuid.UUID) error
		LikeDislikeEpisodeReview(ctx context.Context, arg repository.LikeDislikeEpisodeReviewParams) error
		GetReviewRating(ctx context.Context, arg repository.GetReviewRatingParams) (int64, error)
		IsUserRatedReview(ctx context.Context, arg repository.IsUserRatedReviewParams) (bool, error)

		// Show claps
		AddClapForShow(ctx context.Context, arg repository.AddClapForShowParams) error
		CountUserClaps(ctx context.Context, arg repository.CountUserClapsParams) (int64, error)
	}

	// Challenges service client
	challengesClient interface {
		GetListByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
		NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error)
		GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error)
		GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error)
		ListIDsAvailableUserEpisodes(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]uuid.UUID, error)
	}

	profileClient interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (*profile.Profile, error)
	}

	authClient interface {
		GetUsernameByID(ctx context.Context, uid uuid.UUID) (string, error)
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(sr showsRepository, chc challengesClient, pc profileClient, ac authClient) *Service {
	if sr == nil {
		log.Fatalln("shows repository is not set")
	}
	if chc == nil {
		log.Fatalln("challenges client is not set")
	}
	if pc == nil {
		log.Fatalln("profile client is not set")
	}
	if ac == nil {
		log.Fatalln("auth client is not set")
	}

	return &Service{sr: sr, chc: chc, pc: pc, ac: ac}
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
		sh := Show{
			ID:             s.ID,
			Title:          s.Title,
			Cover:          s.Cover,
			HasNewEpisode:  s.HasNewEpisode,
			Category:       s.Category.String,
			Description:    s.Description.String,
			RealmsTitle:    s.RealmsTitle.String,
			RealmsSubtitle: s.RealmsSubtitle.String,
			Watch:          s.Watch.String,
		}

		if !s.RealmsTitle.Valid {
			sh.RealmsTitle = "Realms"
		}

		result = append(result, sh)
	}

	return result
}

// GetShowByID returns show with provided id.
func (s *Service) GetShowByID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	show, err := s.sr.GetShowByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("could not get show with id=%s: %w", id, err)
	}
	return castToShowWithClaps(show), nil
}

// Cast repository.Show to service Show structure
func castToShow(source repository.Show) Show {
	result := Show{
		ID:             source.ID,
		Title:          source.Title,
		Cover:          source.Cover,
		HasNewEpisode:  source.HasNewEpisode,
		Category:       source.Category.String,
		Description:    source.Description.String,
		RealmsTitle:    source.RealmsTitle.String,
		RealmsSubtitle: source.RealmsSubtitle.String,
		Watch:          source.Watch.String,
	}

	if !source.RealmsTitle.Valid {
		result.RealmsTitle = "Realms"
	}

	return result
}

// Cast repository.GetShowByIDRow to service Show structure
func castToShowWithClaps(source repository.GetShowByIDRow) Show {
	result := Show{
		ID:             source.ID,
		Title:          source.Title,
		Cover:          source.Cover,
		HasNewEpisode:  source.HasNewEpisode,
		Category:       source.Category.String,
		Description:    source.Description.String,
		Claps:          source.Claps,
		RealmsTitle:    source.RealmsTitle.String,
		RealmsSubtitle: source.RealmsSubtitle.String,
		Watch:          source.Watch.String,
	}

	if !source.RealmsTitle.Valid {
		result.RealmsTitle = "Realms"
	}

	return result
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
func (s *Service) GetEpisodesByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error) {
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
		number, err := s.chc.NumberUsersWhoHaveAccessToEpisode(ctx, e.ID)
		if err != nil {
			return nil, fmt.Errorf("could not number users who have access to episode with id = %v: %w", e.ID, err)
		}

		receivedAmount, err := s.chc.GetChallengeReceivedRewardAmount(ctx, e.ChallengeID.UUID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge received reward amount for episode with id = %v: %w", e.ID, err)
		}

		receivedAmountByUser, err := s.chc.GetChallengeReceivedRewardAmountByUserID(ctx, e.ChallengeID.UUID, userID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenge received reward amount by user for episode with id = %v: %w", e.ID, err)
		}

		if _, ok := episodesPerSeasons[e.SeasonID.UUID.String()]; ok {
			episodesPerSeasons[e.SeasonID.UUID.String()] = append(episodesPerSeasons[e.SeasonID.UUID.String()], castRowsToEpisode(e, number, receivedAmount, receivedAmountByUser))
		} else {
			episodesPerSeasons[e.SeasonID.UUID.String()] = []Episode{castRowsToEpisode(e, number, receivedAmount, receivedAmountByUser)}
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
func (s *Service) GetEpisodeByID(ctx context.Context, showID, episodeID, userID uuid.UUID) (Episode, error) {
	episode, err := s.sr.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get episode with id=%s: %w", episodeID, err)
	}

	avgRating, ratingsCount, err := s.getAverageEpisodesRatingByID(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get avarage episoderating with id=%s: %w", episodeID, err)
	}

	receivedAmount, err := s.chc.GetChallengeReceivedRewardAmount(ctx, episode.ChallengeID.UUID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get challenge received reward amount for episode with id = %v: %w", episode.ID, err)
	}

	receivedAmountByUser, err := s.chc.GetChallengeReceivedRewardAmountByUserID(ctx, episode.ChallengeID.UUID, userID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get challenge received reward amount for episode with id = %v: %w", episodeID, err)
	}

	number, err := s.chc.NumberUsersWhoHaveAccessToEpisode(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get number users who have access to episode with id = %v: %w", episodeID, err)
	}

	return castRowToEpisodeExtended(episode, avgRating, receivedAmount, receivedAmountByUser, ratingsCount, number), nil
}

// Cast repository.GetEpisodeByIDRow to service Episode structure with extra data
func castRowToEpisodeExtended(source repository.GetEpisodeByIDRow, rating, receivedAmount, receivedRewardAmountByUser float64, ratingsCount int64, number int32) Episode {
	ep := castRowToEpisode(source)
	ep.Rating = rating
	ep.RatingsCount = ratingsCount
	ep.ActiveUsers = number
	ep.TotalRewardsAmount = receivedAmount
	ep.UserRewardsAmount = receivedRewardAmountByUser

	return ep
}

// Cast repository.GetEpisodeByIDRow to service Episode structure
func castRowToEpisode(source repository.GetEpisodeByIDRow) Episode {
	ep := Episode{
		ID:            source.ID,
		ShowID:        source.ShowID,
		EpisodeNumber: source.EpisodeNumber,
		SeasonID:      source.SeasonID.UUID,
		SeasonNumber:  source.SeasonNumber,
		Cover:         source.Cover.String,
		Title:         source.Title,
		Description:   source.Description.String,
		ReleaseDate:   source.ReleaseDate.Time.String(),
		HintText:      defaultHintText,
		Watch:         source.Watch.String,
	}

	if source.HintText.Valid {
		ep.HintText = source.HintText.String
	}

	if source.ChallengeID.Valid && source.ChallengeID.UUID != uuid.Nil {
		ep.ChallengeID = &source.ChallengeID.UUID
	}

	if source.VerificationChallengeID.Valid && source.VerificationChallengeID.UUID != uuid.Nil {
		ep.VerificationChallengeID = &source.VerificationChallengeID.UUID
	}

	return ep
}

// GetListEpisodesByIDs returns list episodes by list episode ids.
func (s *Service) GetListEpisodesByIDs(ctx context.Context, episodeIDs []uuid.UUID) ([]Episode, error) {
	episodes, err := s.sr.GetListEpisodesByIDs(ctx, episodeIDs)
	if err != nil {
		return []Episode{}, fmt.Errorf("could not get episodes by episodes ids: %w", err)
	}

	result := make([]Episode, 0, len(episodes))

	for _, ep := range episodes {
		result = append(result, Episode{
			ID:                      ep.ID,
			ShowID:                  ep.ShowID,
			SeasonID:                ep.SeasonID.UUID,
			SeasonNumber:            ep.SeasonNumber,
			EpisodeNumber:           ep.EpisodeNumber,
			Cover:                   ep.Cover.String,
			Title:                   ep.Title,
			ShowTitle:               ep.ShowTitle,
			Description:             ep.Description.String,
			ReleaseDate:             ep.ReleaseDate.Time.String(),
			ChallengeID:             &ep.ChallengeID.UUID,
			VerificationChallengeID: &ep.VerificationChallengeID.UUID,
		})
	}

	return result, err
}

// Cast repository.GetEpisodesByShowIDRow to service Episode structure
func castRowsToEpisode(source repository.GetEpisodesByShowIDRow, numberUsersWhoHaveAccessToEpisode int32, receivedAmount, receivedAmountByUser float64) Episode {
	ep := Episode{
		ID:                 source.ID,
		ShowID:             source.ShowID,
		EpisodeNumber:      source.EpisodeNumber,
		SeasonID:           source.SeasonID.UUID,
		SeasonNumber:       source.SeasonNumber,
		Cover:              source.Cover.String,
		Title:              source.Title,
		Description:        source.Description.String,
		ReleaseDate:        source.ReleaseDate.Time.String(),
		Rating:             source.AvgRating,
		RatingsCount:       source.Ratings,
		ActiveUsers:        numberUsersWhoHaveAccessToEpisode,
		TotalRewardsAmount: receivedAmount,
		UserRewardsAmount:  receivedAmountByUser,
		HintText:           defaultHintText,
		Watch:              source.Watch.String,
	}

	if source.HintText.Valid {
		ep.HintText = source.HintText.String
	}

	if source.ChallengeID.Valid && source.ChallengeID.UUID != uuid.Nil {
		ep.ChallengeID = &source.ChallengeID.UUID
	}

	if source.VerificationChallengeID.Valid && source.VerificationChallengeID.UUID != uuid.Nil {
		ep.VerificationChallengeID = &source.VerificationChallengeID.UUID
	}

	return ep
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
		RealmsTitle: sql.NullString{
			String: sh.RealmsTitle,
			Valid:  len(sh.RealmsTitle) > 0,
		},
		RealmsSubtitle: sql.NullString{
			String: sh.RealmsSubtitle,
			Valid:  len(sh.RealmsSubtitle) > 0,
		},
		Watch: sql.NullString{
			String: sh.Watch,
			Valid:  len(sh.Watch) > 0,
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
		RealmsTitle: sql.NullString{
			String: sh.RealmsTitle,
			Valid:  len(sh.RealmsTitle) > 0,
		},
		RealmsSubtitle: sql.NullString{
			String: sh.RealmsSubtitle,
			Valid:  len(sh.RealmsSubtitle) > 0,
		},
		Watch: sql.NullString{
			String: sh.Watch,
			Valid:  len(sh.Watch) > 0,
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

	params := repository.AddEpisodeParams{
		ShowID: ep.ShowID,
		SeasonID: uuid.NullUUID{
			UUID:  ep.SeasonID,
			Valid: ep.SeasonID != uuid.Nil,
		},
		EpisodeNumber: ep.EpisodeNumber,
		Cover: sql.NullString{
			String: ep.Cover,
			Valid:  len(ep.Cover) > 0,
		},
		Title: ep.Title,
		Description: sql.NullString{
			String: ep.Description,
			Valid:  len(ep.Description) > 0,
		},
		ReleaseDate: sql.NullTime{
			Time:  rDate,
			Valid: true,
		},
		HintText: sql.NullString{
			String: ep.HintText,
			Valid:  len(ep.HintText) > 0 && ep.HintText != defaultHintText,
		},
		Watch: sql.NullString{
			String: ep.Watch,
			Valid:  len(ep.Watch) > 0,
		},
	}

	if ep.ChallengeID != nil && *ep.ChallengeID != uuid.Nil {
		params.ChallengeID = uuid.NullUUID{UUID: *ep.ChallengeID, Valid: true}
	}

	if ep.VerificationChallengeID != nil && *ep.VerificationChallengeID != uuid.Nil {
		params.VerificationChallengeID = uuid.NullUUID{UUID: *ep.VerificationChallengeID, Valid: true}
	}

	episode, err := s.sr.AddEpisode(ctx, params)
	if err != nil {
		return Episode{}, fmt.Errorf("could not add episode #%d for show_id=%s: %w", ep.EpisodeNumber, ep.ShowID.String(), err)
	}

	episodeByID, err := s.sr.GetEpisodeByID(ctx, episode.ID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get episode with id=%s: %w", episode.ID, err)
	}

	// return castToEpisode(episode, episodeByID.SeasonNumber), nil
	return castRowToEpisode(episodeByID), nil
}

// UpdateEpisode ..
func (s *Service) UpdateEpisode(ctx context.Context, ep Episode) error {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return fmt.Errorf("could not add parse date from string: %w", err)
	}

	params := repository.UpdateEpisodeParams{
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
		HintText: sql.NullString{
			String: ep.HintText,
			Valid:  len(ep.HintText) > 0 && ep.HintText != defaultHintText,
		},
		Watch: sql.NullString{
			String: ep.Watch,
			Valid:  len(ep.Watch) > 0,
		},
	}

	if ep.ChallengeID != nil && *ep.ChallengeID != uuid.Nil {
		params.ChallengeID = uuid.NullUUID{UUID: *ep.ChallengeID, Valid: true}
	}

	if ep.VerificationChallengeID != nil && *ep.VerificationChallengeID != uuid.Nil {
		params.VerificationChallengeID = uuid.NullUUID{UUID: *ep.VerificationChallengeID, Valid: true}
	}

	if err = s.sr.UpdateEpisode(ctx, params); err != nil {
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
			return fmt.Errorf("you've already rated this episode")
		}
		return fmt.Errorf("could not rate episode with episodeID=%s: %w", episodeID, err)
	}

	return nil
}

// ReviewEpisode ...
func (s *Service) ReviewEpisode(ctx context.Context, episodeID, userID uuid.UUID, username string, rating int32, title, review string) error {
	if err := s.sr.ReviewEpisode(ctx, repository.ReviewEpisodeParams{
		EpisodeID: episodeID,
		UserID:    userID,
		Username:  sql.NullString{String: username, Valid: true},
		Rating:    rating,
		Title:     sql.NullString{String: title, Valid: true},
		Review:    sql.NullString{String: review, Valid: true},
	}); err != nil {
		return fmt.Errorf("could not review episode with episodeID=%s: %w", episodeID, err)
	}

	return nil
}

func (s *Service) GetReviewsList(ctx context.Context, episodeID uuid.UUID, limit, offset int32, currentUserID uuid.UUID) ([]Review, error) {
	reviews, err := s.sr.ReviewsList(ctx, repository.ReviewsListParams{
		EpisodeID: episodeID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.castReviewsList(ctx, reviews, currentUserID), nil
}

func (s *Service) castReviewsList(ctx context.Context, source []repository.Rating, currentUserID uuid.UUID) []Review {
	result := make([]Review, 0, len(source))
	for _, r := range source {
		result = append(result, s.castReview(ctx, r, currentUserID))
	}

	return result
}

func (s *Service) castReview(ctx context.Context, source repository.Rating, currentUserID uuid.UUID) Review {
	profile, err := s.pc.GetProfileByUserID(ctx, source.UserID, "")
	if err != nil {
		log.Printf("could not get profile by user id: %v", err)
	}
	username, err := s.ac.GetUsernameByID(ctx, source.UserID)
	if err != nil {
		log.Printf("could not get username by user id: %v", err)
	}
	likes, _ := s.sr.GetReviewRating(ctx, repository.GetReviewRatingParams{
		ReviewID: source.ID,
		RatingType: sql.NullInt32{
			Int32: int32(LikeReview),
			Valid: true,
		},
	})
	dislikes, _ := s.sr.GetReviewRating(ctx, repository.GetReviewRatingParams{
		ReviewID: source.ID,
		RatingType: sql.NullInt32{
			Int32: int32(DislikeReview),
			Valid: true,
		},
	})
	isLiked, _ := s.sr.IsUserRatedReview(ctx, repository.IsUserRatedReviewParams{
		UserID:   currentUserID,
		ReviewID: source.ID,
		RatingType: sql.NullInt32{
			Int32: int32(LikeReview),
			Valid: true,
		},
	})
	isDisliked, _ := s.sr.IsUserRatedReview(ctx, repository.IsUserRatedReviewParams{
		UserID:   currentUserID,
		ReviewID: source.ID,
		RatingType: sql.NullInt32{
			Int32: int32(DislikeReview),
			Valid: true,
		},
	})
	return Review{
		ID:         source.ID.String(),
		UserID:     source.UserID.String(),
		Username:   username,
		UserAvatar: profile.Avatar,
		Rating:     int(source.Rating),
		Title:      source.Title.String,
		Review:     source.Review.String,
		CreatedAt:  source.CreatedAt.Format(time.RFC3339),
		Likes:      likes,
		Dislikes:   dislikes,
		IsLiked:    isLiked,
		IsDisliked: isDisliked,
	}
}

// DeleteReviewByID ..
func (s *Service) DeleteReviewByID(ctx context.Context, id uuid.UUID) error {
	if err := s.sr.DeleteReview(ctx, id); err != nil {
		return fmt.Errorf("could not delete review with id=%s:%w", id, err)
	}

	return nil
}

// AddClapsForShow ...
func (s *Service) AddClapsForShow(ctx context.Context, showID, userID uuid.UUID) error {
	if claps, err := s.sr.CountUserClaps(ctx, repository.CountUserClapsParams{
		ShowID: showID,
		UserID: userID,
	}); err == nil {
		if claps >= 10 {
			return ErrMaxClaps
		}
	}

	if err := s.sr.AddClapForShow(ctx, repository.AddClapForShowParams{
		ShowID: showID,
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("could not add claps for show with id=%s: %w", showID, err)
	}

	return nil
}

// GetReviewsListByUserID ...
func (s *Service) GetReviewsListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32, currentUserID uuid.UUID) ([]Review, error) {
	reviews, err := s.sr.ReviewsListByUserID(ctx, repository.ReviewsListByUserIDParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return s.castReviewsList(ctx, reviews, currentUserID), nil
}

// LikeDislikeEpisodeReview used to store users review episode assessment (like/dislike).
func (s *Service) LikeDislikeEpisodeReview(ctx context.Context, reviewID, uid uuid.UUID, ratingType ReviewRatingType) error {
	err := s.sr.LikeDislikeEpisodeReview(ctx, repository.LikeDislikeEpisodeReviewParams{
		ReviewID: reviewID,
		UserID:   uid,
		RatingType: sql.NullInt32{
			Int32: int32(ratingType),
			Valid: true,
		}})
	if err != nil {
		return fmt.Errorf("could not like/dislike review episode with id=:%v, %w", uid, err)
	}

	return nil
}

// GetActivatedUserEpisodes returns list activated episodes by user id.
func (s *Service) GetActivatedUserEpisodes(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Episode, error) {
	listIDs, err := s.chc.ListIDsAvailableUserEpisodes(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("could not get list ids user available episodes: %w", err)
	}

	listEpisodes, err := s.GetListEpisodesByIDs(ctx, listIDs)
	if err != nil {
		return nil, fmt.Errorf("could not get list user available episodes: %w", err)
	}

	return listEpisodes, nil
}
