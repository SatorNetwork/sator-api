package shows

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/db"
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/SatorNetwork/sator-api/svc/profile"
	"github.com/SatorNetwork/sator-api/svc/shows/repository"
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
		sr           showsRepository
		chc          challengesClient
		pc           profileClient
		ac           authClient
		sentTipsFunc sentTipsFunction
		nc           nftClient
		firebaseSvc  firebaseService
		tipsPercent  float64
	}

	// Show struct
	// Fields were rearranged to optimize memory usage.
	Show struct {
		ID             uuid.UUID   `json:"id"`
		Title          string      `json:"title"`
		Cover          string      `json:"cover"`
		HasNewEpisode  bool        `json:"has_new_episode"`
		Categories     []uuid.UUID `json:"categories"`
		Description    string      `json:"description"`
		Claps          int64       `json:"claps"`
		RealmsTitle    string      `json:"realms_title"`
		RealmsSubtitle string      `json:"realms_subtitle"`
		Watch          string      `json:"watch"`
		HasNFT         bool        `json:"has_nft"`
		Status         string      `json:"status"`
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
		UsersEpisodeRating      int32      `json:"users_episode_rating"`
		ActiveUsers             int32      `json:"active_users"`
		UserRewardsAmount       float64    `json:"user_rewards_amount"`
		TotalRewardsAmount      float64    `json:"total_rewards_amount"`
		HintText                string     `json:"hint_text"`
		Watch                   string     `json:"watch"`
		Status                  string     `json:"status"`
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

	ShowCategory struct {
		ID       uuid.UUID `json:"id"`
		Title    string    `json:"title"`
		Disabled bool      `json:"disabled"`
		Sort     int32     `json:"sort"`
	}

	showsRepository interface {
		// Shows
		AddShow(ctx context.Context, arg repository.AddShowParams) (repository.Show, error)
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetAllShows(ctx context.Context, arg repository.GetAllShowsParams) ([]repository.Show, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.Show, error)
		GetPublishedShows(ctx context.Context, arg repository.GetPublishedShowsParams) ([]repository.Show, error)
		GetPublishedShowByID(ctx context.Context, id uuid.UUID) (repository.GetPublishedShowByIDRow, error)
		GetShowsByCategory(ctx context.Context, arg repository.GetShowsByCategoryParams) ([]repository.Show, error)
		UpdateShow(ctx context.Context, arg repository.UpdateShowParams) error
		GetShowsByOldCategory(ctx context.Context, arg repository.GetShowsByOldCategoryParams) ([]repository.Show, error)
		GetShowsByStatus(ctx context.Context, arg repository.GetShowsByStatusParams) ([]repository.Show, error)

		// Seasons
		AddSeason(ctx context.Context, arg repository.AddSeasonParams) (repository.Season, error)
		DeleteSeasonByID(ctx context.Context, id uuid.UUID) error
		DeleteSeasonByShowID(ctx context.Context, showID uuid.UUID) error
		GetSeasonByID(ctx context.Context, id uuid.UUID) (repository.Season, error)
		GetSeasonsByShowID(ctx context.Context, arg repository.GetSeasonsByShowIDParams) ([]repository.Season, error)

		// Episodes
		AddEpisode(ctx context.Context, arg repository.AddEpisodeParams) (repository.Episode, error)
		GetEpisodeByID(ctx context.Context, id uuid.UUID) (repository.GetEpisodeByIDRow, error)
		GetPublishedEpisodeByID(ctx context.Context, id uuid.UUID) (repository.GetPublishedEpisodeByIDRow, error)
		GetPublishedListEpisodesByIDs(ctx context.Context, episodeIds []uuid.UUID) ([]repository.GetPublishedListEpisodesByIDsRow, error)
		GetPublishedEpisodesByShowID(ctx context.Context, arg repository.GetPublishedEpisodesByShowIDParams) ([]repository.GetPublishedEpisodesByShowIDRow, error)
		GetAllEpisodesByShowID(ctx context.Context, arg repository.GetAllEpisodesByShowIDParams) ([]repository.GetAllEpisodesByShowIDRow, error)
		GetEpisodesByStatus(ctx context.Context, arg repository.GetEpisodesByStatusParams) ([]repository.Episode, error)
		DeleteEpisodeByID(ctx context.Context, id uuid.UUID) error
		DeleteEpisodeByShowID(ctx context.Context, showID uuid.UUID) error
		DeleteEpisodeBySeasonID(ctx context.Context, seasonID uuid.NullUUID) error
		UpdateEpisode(ctx context.Context, arg repository.UpdateEpisodeParams) error
		LinkEpisodeChallenges(ctx context.Context, arg repository.LinkEpisodeChallengesParams) error

		// Episodes rating
		GetEpisodeRatingByID(ctx context.Context, episodeID uuid.UUID) (repository.GetEpisodeRatingByIDRow, error)
		RateEpisode(ctx context.Context, arg repository.RateEpisodeParams) error
		DidUserRateEpisode(ctx context.Context, arg repository.DidUserRateEpisodeParams) (bool, error)
		GetUsersEpisodeRatingByID(ctx context.Context, arg repository.GetUsersEpisodeRatingByIDParams) (int32, error)

		// Episode reviews
		DidUserReviewEpisode(ctx context.Context, arg repository.DidUserReviewEpisodeParams) (bool, error)
		ReviewEpisode(ctx context.Context, arg repository.ReviewEpisodeParams) (repository.Rating, error)
		ReviewsList(ctx context.Context, arg repository.ReviewsListParams) ([]repository.ReviewsListRow, error)
		ReviewsListByUserID(ctx context.Context, arg repository.ReviewsListByUserIDParams) ([]repository.ReviewsListByUserIDRow, error)
		DeleteReview(ctx context.Context, id uuid.UUID) error
		LikeDislikeEpisodeReview(ctx context.Context, arg repository.LikeDislikeEpisodeReviewParams) error
		GetReviewRating(ctx context.Context, arg repository.GetReviewRatingParams) (int64, error)
		IsUserRatedReview(ctx context.Context, arg repository.IsUserRatedReviewParams) (bool, error)
		GetReviewByID(ctx context.Context, id uuid.UUID) (repository.Rating, error)
		GetUserEpisodeReview(ctx context.Context, arg repository.GetUserEpisodeReviewParams) (repository.ReviewsRating, error)
		DeleteUserEpisodeReview(ctx context.Context, arg repository.DeleteUserEpisodeReviewParams) error

		// Show claps
		AddClapForShow(ctx context.Context, arg repository.AddClapForShowParams) error
		CountUserClaps(ctx context.Context, arg repository.CountUserClapsParams) (int64, error)

		// Show category
		AddShowCategory(ctx context.Context, arg repository.AddShowCategoryParams) (repository.ShowCategory, error)
		DeleteShowCategoryByID(ctx context.Context, id uuid.UUID) error
		GetShowCategories(ctx context.Context, arg repository.GetShowCategoriesParams) ([]repository.ShowCategory, error)
		GetShowCategoriesWithDisabled(ctx context.Context, arg repository.GetShowCategoriesWithDisabledParams) ([]repository.ShowCategory, error)
		GetShowCategoryByID(ctx context.Context, id uuid.UUID) (repository.ShowCategory, error)
		UpdateShowCategory(ctx context.Context, arg repository.UpdateShowCategoryParams) error

		//	Show to category
		AddShowToCategory(ctx context.Context, arg repository.AddShowToCategoryParams) (repository.ShowsToCategory, error)
		DeleteShowToCategoryByShowID(ctx context.Context, showID uuid.UUID) error
		GetCategoriesByShowID(ctx context.Context, showID uuid.UUID) ([]uuid.UUID, error)
	}

	// Challenges service client
	challengesClient interface {
		GetListByShowID(ctx context.Context, showID, userID uuid.UUID, limit, offset int32) (interface{}, error)
		NumberUsersWhoHaveAccessToEpisode(ctx context.Context, episodeID uuid.UUID) (int32, error)
		GetChallengeReceivedRewardAmount(ctx context.Context, challengeID uuid.UUID) (float64, error)
		GetChallengeReceivedRewardAmountByUserID(ctx context.Context, challengeID, userID uuid.UUID) (float64, error)
		ListIDsAvailableUserEpisodes(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]uuid.UUID, error)

		CreateQuizChallenge(ctx context.Context, showID, epID uuid.UUID, episodeTitle string) (uuid.UUID, error)
		CreateVerificationChallenge(ctx context.Context, showID, epID uuid.UUID, episodeTitle string) (uuid.UUID, error)
	}

	profileClient interface {
		GetProfileByUserID(ctx context.Context, userID uuid.UUID, username string) (*profile.Profile, error)
	}

	authClient interface {
		GetUsernameByID(ctx context.Context, uid uuid.UUID) (string, error)
	}

	// NFT service client
	nftClient interface {
		DoesRelationIDHasNFT(ctx context.Context, relationID uuid.UUID) (bool, error)
	}

	firebaseService interface {
		SendNewShowNotification(ctx context.Context, showTitle string, showID uuid.UUID) error
		SendNewEpisodeNotification(ctx context.Context, showTitle, episodeTitle string, showID, seasonID, episodeID uuid.UUID) error
	}

	// Simple function
	sentTipsFunction func(ctx context.Context, uid, recipientID uuid.UUID, amount float64, cfg *lib_solana.SendAssetsConfig, info string) error
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation.
func NewService(
	sr showsRepository,
	chc challengesClient,
	pc profileClient,
	ac authClient,
	sentTipsFunc sentTipsFunction,
	nc nftClient,
	firebaseSvc firebaseService,
	tipsPercent float64,
) *Service {
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
	if sentTipsFunc == nil {
		log.Fatalln("sentTipsFunc is not set")
	}
	if nc == nil {
		log.Fatalln("nft client is not set")
	}

	return &Service{
		sr:           sr,
		chc:          chc,
		pc:           pc,
		ac:           ac,
		sentTipsFunc: sentTipsFunc,
		nc:           nc,
		tipsPercent:  tipsPercent,
		firebaseSvc:  firebaseSvc,
	}
}

// GetShows returns shows.
func (s *Service) GetShows(ctx context.Context, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetPublishedShows(ctx, repository.GetPublishedShowsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	sl, err := s.castToListShow(ctx, shows)
	if err != nil {
		return nil, fmt.Errorf("could not cast to list show : %w", err)
	}

	return sl, nil
}

// GetAllShows returns shows.
func (s *Service) GetAllShows(ctx context.Context, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetAllShows(ctx, repository.GetAllShowsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	sl, err := s.castToListShow(ctx, shows)
	if err != nil {
		return nil, fmt.Errorf("could not cast to list show : %w", err)
	}

	return sl, nil
}

// GetShowsWithNFT returns shows list which has NFT.
func (s *Service) GetShowsWithNFT(ctx context.Context, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetPublishedShows(ctx, repository.GetPublishedShowsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	result := make([]Show, 0, len(shows))
	for _, sw := range shows {
		hasNFT, err := s.nc.DoesRelationIDHasNFT(ctx, sw.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenges list by show id: %v", err)
		}
		if !hasNFT {
			continue
		}
		sh := Show{
			ID:             sw.ID,
			Title:          sw.Title,
			Cover:          sw.Cover,
			HasNewEpisode:  sw.HasNewEpisode,
			Description:    sw.Description.String,
			RealmsTitle:    sw.RealmsTitle.String,
			RealmsSubtitle: sw.RealmsSubtitle.String,
			Watch:          sw.Watch.String,
			HasNFT:         hasNFT,
			Status:         string(sw.Status),
		}

		if !sw.RealmsTitle.Valid {
			sh.RealmsTitle = "Realms"
		}

		result = append(result, sh)
	}

	return result, nil
}

func (s *Service) GetShowsByStatus(ctx context.Context, status string, limit, offset int32) ([]Show, error) {
	shows, err := s.sr.GetShowsByStatus(ctx, repository.GetShowsByStatusParams{
		Status: repository.ShowsStatusType(status),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	result, err := s.castToListShow(ctx, shows)
	if err != nil {
		return nil, fmt.Errorf("could not cast to list show : %w", err)
	}

	return result, nil
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
func (s *Service) castToListShow(ctx context.Context, source []repository.Show) ([]Show, error) {
	result := make([]Show, 0, len(source))
	for _, sw := range source {
		hasNFT, err := s.nc.DoesRelationIDHasNFT(ctx, sw.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get challenges list by show id: %v", err)
		}
		sh := Show{
			ID:             sw.ID,
			Title:          sw.Title,
			Cover:          sw.Cover,
			HasNewEpisode:  sw.HasNewEpisode,
			Description:    sw.Description.String,
			RealmsTitle:    sw.RealmsTitle.String,
			RealmsSubtitle: sw.RealmsSubtitle.String,
			Watch:          sw.Watch.String,
			HasNFT:         hasNFT,
			Status:         string(sw.Status),
		}

		if !sw.RealmsTitle.Valid {
			sh.RealmsTitle = "Realms"
		}

		result = append(result, sh)
	}

	return result, nil
}

// GetShowByID returns show with provided id.
func (s *Service) GetShowByID(ctx context.Context, id uuid.UUID) (Show, error) {
	show, err := s.sr.GetShowByID(ctx, id)
	if err != nil {
		return Show{}, fmt.Errorf("could not get show with id=%s: %w", id, err)
	}

	result := Show{
		ID:             show.ID,
		Title:          show.Title,
		Cover:          show.Cover,
		HasNewEpisode:  show.HasNewEpisode,
		Description:    show.Description.String,
		RealmsTitle:    show.RealmsTitle.String,
		RealmsSubtitle: show.RealmsSubtitle.String,
		Watch:          show.Watch.String,
		Status:         string(show.Status),
	}

	if !show.RealmsTitle.Valid {
		result.RealmsTitle = "Realms"
	}

	categories, err := s.sr.GetCategoriesByShowID(ctx, id)
	if err != nil {
		return Show{}, fmt.Errorf("could not get categories list by show id: %v", err)
	}

	for i := 0; i < len(categories); i++ {
		result.Categories = append(result.Categories, categories[i])
	}

	return result, nil
}

// GetPublishedShowByID returns show with provided id.
func (s *Service) GetPublishedShowByID(ctx context.Context, id uuid.UUID) (Show, error) {
	show, err := s.sr.GetPublishedShowByID(ctx, id)
	if err != nil {
		return Show{}, fmt.Errorf("could not get show with id=%s: %w", id, err)
	}
	hasNFT, err := s.nc.DoesRelationIDHasNFT(ctx, show.ID)
	if err != nil {
		return Show{}, fmt.Errorf("could not get challenges list by show id: %v", err)
	}

	result := Show{
		ID:             show.ID,
		Title:          show.Title,
		Cover:          show.Cover,
		HasNewEpisode:  show.HasNewEpisode,
		Description:    show.Description.String,
		Claps:          show.Claps,
		RealmsTitle:    show.RealmsTitle.String,
		RealmsSubtitle: show.RealmsSubtitle.String,
		Watch:          show.Watch.String,
		HasNFT:         hasNFT,
		Status:         string(show.Status),
	}

	if !show.RealmsTitle.Valid {
		result.RealmsTitle = "Realms"
	}

	categories, err := s.sr.GetCategoriesByShowID(ctx, id)
	if err != nil {
		return Show{}, fmt.Errorf("could not get categories list by show id: %v", err)
	}

	for i := 0; i < len(categories); i++ {
		result.Categories = append(result.Categories, categories[i])
	}

	return result, nil
}

// GetShowsByCategory returns show by provided category.
func (s *Service) GetShowsByCategory(ctx context.Context, category uuid.UUID, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetShowsByCategory(ctx, repository.GetShowsByCategoryParams{
		CategoryID: category,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	sl, err := s.castToListShow(ctx, shows)
	if err != nil {
		return nil, fmt.Errorf("could not cast to list show : %w", err)
	}

	return sl, nil
}

// GetShowsByOldCategory returns show by provided category.
// TODO: DEPRECATED, will be removed in one of the following releases
func (s *Service) GetShowsByOldCategory(ctx context.Context, category string, limit, offset int32) (interface{}, error) {
	shows, err := s.sr.GetShowsByOldCategory(ctx, repository.GetShowsByOldCategoryParams{
		Category: category,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	sl, err := s.castToListShow(ctx, shows)
	if err != nil {
		return nil, fmt.Errorf("could not cast to list show : %w", err)
	}

	return sl, nil
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

	episodes, err := s.sr.GetPublishedEpisodesByShowID(ctx, repository.GetPublishedEpisodesByShowIDParams{
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

		ep := castRowsToPublishedEpisode(e, number, receivedAmount, receivedAmountByUser)
		if _, ok := episodesPerSeasons[e.SeasonID.UUID.String()]; ok {
			episodesPerSeasons[e.SeasonID.UUID.String()] = append(episodesPerSeasons[e.SeasonID.UUID.String()], ep)
		} else {
			episodesPerSeasons[e.SeasonID.UUID.String()] = []Episode{ep}
		}
	}

	return castToListSeasons(seasons, episodesPerSeasons), nil
}

// GetAllEpisodesByShowID returns episodes by show id.
func (s *Service) GetAllEpisodesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	seasons, err := s.sr.GetSeasonsByShowID(ctx, repository.GetSeasonsByShowIDParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get seasons list: %w", err)
	}

	episodes, err := s.sr.GetAllEpisodesByShowID(ctx, repository.GetAllEpisodesByShowIDParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get episodes list: %w", err)
	}

	episodesPerSeasons := make(map[string][]Episode)
	for _, e := range episodes {
		ep := castRowsToEpisode(e)
		if _, ok := episodesPerSeasons[e.SeasonID.UUID.String()]; ok {
			episodesPerSeasons[e.SeasonID.UUID.String()] = append(episodesPerSeasons[e.SeasonID.UUID.String()], ep)
		} else {
			episodesPerSeasons[e.SeasonID.UUID.String()] = []Episode{ep}
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

func (s *Service) GetEpisodesByStatus(ctx context.Context, status string, limit, offset int32) ([]Episode, error) {
	episodes, err := s.sr.GetEpisodesByStatus(ctx, repository.GetEpisodesByStatusParams{
		Status: repository.EpisodesStatusType(status),
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("could not get shows list: %w", err)
	}

	result := make([]Episode, 0, len(episodes))
	for _, ep := range episodes {
		result = append(result, Episode{
			ID:                      ep.ID,
			ShowID:                  ep.ShowID,
			SeasonID:                ep.SeasonID.UUID,
			EpisodeNumber:           ep.EpisodeNumber,
			Cover:                   ep.Cover.String,
			Title:                   ep.Title,
			Description:             ep.Description.String,
			ReleaseDate:             ep.ReleaseDate.Time.String(),
			ChallengeID:             &ep.ChallengeID.UUID,
			VerificationChallengeID: &ep.VerificationChallengeID.UUID,
		})
	}

	return result, nil
}

// GetEpisodeByID returns episode with provided id.
func (s *Service) GetEpisodeByID(ctx context.Context, episodeID, userID uuid.UUID) (Episode, error) {
	episode, err := s.sr.GetEpisodeByID(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get episode with id=%s: %w", episodeID, err)
	}

	needsToUpdate := false

	if !episode.ChallengeID.Valid || episode.ChallengeID.UUID == uuid.Nil {
		quizID, err := s.chc.CreateQuizChallenge(ctx, episode.ShowID, episode.ID, episode.Title)
		if err != nil {
			return Episode{}, fmt.Errorf("could not create quiz challenge for episode with id=%s: %w", episodeID, err)
		}

		episode.ChallengeID = uuid.NullUUID{UUID: quizID, Valid: true}
		needsToUpdate = true
	}

	if !episode.VerificationChallengeID.Valid || episode.VerificationChallengeID.UUID == uuid.Nil {
		verificationID, err := s.chc.CreateVerificationChallenge(ctx, episode.ShowID, episode.ID, episode.Title)
		if err != nil {
			return Episode{}, fmt.Errorf("could not create verification challenge for episode with id=%s: %w", episodeID, err)
		}

		episode.VerificationChallengeID = uuid.NullUUID{UUID: verificationID, Valid: true}
		needsToUpdate = true
	}

	if needsToUpdate {
		if err := s.sr.LinkEpisodeChallenges(ctx, repository.LinkEpisodeChallengesParams{
			ID:                      episode.ID,
			ChallengeID:             episode.ChallengeID,
			VerificationChallengeID: episode.VerificationChallengeID,
		}); err != nil {
			return Episode{}, fmt.Errorf("could not update chalenges ids for episode with id=%s: %w", episodeID, err)
		}
	}

	return castRowToEpisode(episode), nil
}

// GetPublishedEpisodeByID returns episode with provided id.
func (s *Service) GetPublishedEpisodeByID(ctx context.Context, episodeID, userID uuid.UUID) (Episode, error) {
	episode, err := s.sr.GetPublishedEpisodeByID(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get episode with id=%s: %w", episodeID, err)
	}

	needsToUpdate := false

	if !episode.ChallengeID.Valid || episode.ChallengeID.UUID == uuid.Nil {
		quizID, err := s.chc.CreateQuizChallenge(ctx, episode.ShowID, episode.ID, episode.Title)
		if err != nil {
			return Episode{}, fmt.Errorf("could not create quiz challenge for episode with id=%s: %w", episodeID, err)
		}

		episode.ChallengeID = uuid.NullUUID{UUID: quizID, Valid: true}
		needsToUpdate = true
	}

	if !episode.VerificationChallengeID.Valid || episode.VerificationChallengeID.UUID == uuid.Nil {
		verificationID, err := s.chc.CreateVerificationChallenge(ctx, episode.ShowID, episode.ID, episode.Title)
		if err != nil {
			return Episode{}, fmt.Errorf("could not create verification challenge for episode with id=%s: %w", episodeID, err)
		}

		episode.VerificationChallengeID = uuid.NullUUID{UUID: verificationID, Valid: true}
		needsToUpdate = true
	}

	if needsToUpdate {
		if err := s.sr.LinkEpisodeChallenges(ctx, repository.LinkEpisodeChallengesParams{
			ID:                      episode.ID,
			ChallengeID:             episode.ChallengeID,
			VerificationChallengeID: episode.VerificationChallengeID,
		}); err != nil {
			return Episode{}, fmt.Errorf("could not update chalenges ids for episode with id=%s: %w", episodeID, err)
		}
	}

	avgRating, ratingsCount, err := s.getAverageEpisodesRatingByID(ctx, episodeID)
	if err != nil {
		return Episode{}, fmt.Errorf("could not get avarage episode rating with id=%s: %w", episodeID, err)
	}

	usersEpisodeRating, err := s.sr.GetUsersEpisodeRatingByID(ctx, repository.GetUsersEpisodeRatingByIDParams{
		EpisodeID: episodeID,
		UserID:    userID,
	})
	if err != nil && !db.IsNotFoundError(err) {
		return Episode{}, fmt.Errorf("could not get users episode rating with id=%s: %w", episodeID, err)
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

	return castRowToEpisodeExtended(episode, avgRating, receivedAmount, receivedAmountByUser, ratingsCount, number, usersEpisodeRating), nil
}

// Cast repository.GetEpisodeByIDRow to service Episode structure with extra data
func castRowToEpisodeExtended(source repository.GetPublishedEpisodeByIDRow, rating, receivedAmount, receivedRewardAmountByUser float64, ratingsCount int64, number, usersEpisodeRating int32) Episode {
	ep := castPublishedRowToEpisode(source)
	ep.Rating = rating
	ep.RatingsCount = ratingsCount
	ep.ActiveUsers = number
	ep.TotalRewardsAmount = receivedAmount
	ep.UserRewardsAmount = receivedRewardAmountByUser
	ep.UsersEpisodeRating = usersEpisodeRating

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
		Status:        string(source.Status),
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

// Cast repository.GetPublishedEpisodeByIDRow to service Episode structure
func castPublishedRowToEpisode(source repository.GetPublishedEpisodeByIDRow) Episode {
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
		Status:        string(source.Status),
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
	episodes, err := s.sr.GetPublishedListEpisodesByIDs(ctx, episodeIDs)
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
func castRowsToPublishedEpisode(source repository.GetPublishedEpisodesByShowIDRow, numberUsersWhoHaveAccessToEpisode int32, receivedAmount, receivedAmountByUser float64) Episode {
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
		Status:             string(source.Status),
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

// Cast repository.GetEpisodesByShowIDRow to service Episode structure
func castRowsToEpisode(source repository.GetAllEpisodesByShowIDRow) Episode {
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
		Status:        string(source.Status),
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
		Status: repository.ShowsStatusType(sh.Status),
	})
	if err != nil {
		return Show{}, fmt.Errorf("could not add show with title=%s: %w", sh.Title, err)
	}

	for _, cat := range sh.Categories {
		if _, err = s.sr.AddShowToCategory(ctx, repository.AddShowToCategoryParams{
			CategoryID: cat,
			ShowID:     show.ID,
		}); err != nil && !db.IsNotFoundError(err) {
			log.Printf("could not add category to show with show id=%s: %v", show.ID, err)
		}
	}

	if show.Status == repository.ShowsStatusTypePublished {
		if err := s.firebaseSvc.SendNewShowNotification(ctx, show.Title, show.ID); err != nil {
			log.Printf("could not send new show notification: %v", err)
		}
	}

	return s.GetShowByID(ctx, show.ID)
}

// UpdateShow ...
func (s *Service) UpdateShow(ctx context.Context, sh Show) error {
	oldShow, err := s.sr.GetShowByID(ctx, sh.ID)
	if err != nil {
		return fmt.Errorf("could not get show with id=%s: %w", sh.ID, err)
	}

	if err := s.sr.UpdateShow(ctx, repository.UpdateShowParams{
		Title:         sh.Title,
		Cover:         sh.Cover,
		HasNewEpisode: sh.HasNewEpisode,
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
		Status: repository.ShowsStatusType(sh.Status),
		ID:     sh.ID,
	}); err != nil {
		return fmt.Errorf("could not update show with id=%s:%w", sh.ID, err)
	}

	if err := s.sr.DeleteShowToCategoryByShowID(ctx, sh.ID); err != nil && !db.IsNotFoundError(err) {
		return fmt.Errorf("could not delete categories with show id=%s: %w", sh.ID, err)
	}

	for _, cat := range sh.Categories {
		if _, err = s.sr.AddShowToCategory(ctx, repository.AddShowToCategoryParams{
			CategoryID: cat,
			ShowID:     sh.ID,
		}); err != nil && !db.IsNotFoundError(err) {
			log.Printf("could not add category to show with show id=%s: %v", sh.ID, err)
		}
	}

	if oldShow.Status != repository.ShowsStatusTypePublished &&
		repository.ShowsStatusType(sh.Status) == repository.ShowsStatusTypePublished {
		if err := s.firebaseSvc.SendNewShowNotification(ctx, sh.Title, sh.ID); err != nil {
			log.Printf("can't send new show notification: %v", err)
		}
	}

	return nil
}

// DeleteShowByID ..
func (s *Service) DeleteShowByID(ctx context.Context, id uuid.UUID) error {
	if err := s.sr.DeleteShowByID(ctx, id); err != nil {
		return fmt.Errorf("could not delete show with id=%s:%w", id, err)
	}

	if err := s.sr.DeleteSeasonByShowID(ctx, id); err != nil {
		return fmt.Errorf("could not delete seasons with show id=%s:%w", id, err)
	}

	if err := s.sr.DeleteEpisodeByShowID(ctx, id); err != nil {
		return fmt.Errorf("could not delete episodes with show id=%s:%w", id, err)
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
		Status: repository.EpisodesStatusType(ep.Status),
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

	show, err := s.sr.GetShowByID(ctx, ep.ShowID)
	if err != nil {
		return Episode{}, errors.Wrap(err, "can't get show by id")
	}

	if show.Status == repository.ShowsStatusTypePublished && episode.Status == repository.EpisodesStatusTypePublished {
		if err := s.firebaseSvc.SendNewEpisodeNotification(ctx, show.Title, ep.Title, show.ID, ep.SeasonID, episode.ID); err != nil {
			log.Printf("can't send new episode notification: %v", err)
		}
	}

	return castRowToEpisode(episodeByID), nil
}

// UpdateEpisode ..
func (s *Service) UpdateEpisode(ctx context.Context, ep Episode) error {
	rDate, err := utils.DateFromString(ep.ReleaseDate)
	if err != nil {
		return fmt.Errorf("could not add parse date from string: %w", err)
	}

	oldEpisode, err := s.sr.GetEpisodeByID(ctx, ep.ID)
	if err != nil {
		return fmt.Errorf("could not get episode with id=%s:%w", ep.ID, err)
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
		Status: repository.EpisodesStatusType(ep.Status),
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

	show, err := s.sr.GetShowByID(ctx, ep.ShowID)
	if err != nil {
		return errors.Wrap(err, "can't get show by id")
	}

	if show.Status == repository.ShowsStatusTypePublished &&
		oldEpisode.Status != repository.EpisodesStatusTypePublished &&
		repository.EpisodesStatusType(ep.Status) == repository.EpisodesStatusTypePublished {

		if err := s.firebaseSvc.SendNewEpisodeNotification(ctx, show.Title, ep.Title, show.ID, ep.SeasonID, ep.ID); err != nil {
			log.Printf("can't send new show notification: %v", err)
		}
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

func (s *Service) GetSeasonByID(ctx context.Context, showID, seasonID uuid.UUID) (*Season, error) {
	season, err := s.sr.GetSeasonByID(ctx, seasonID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get season by id")
	}
	if showID != season.ShowID {
		return nil, errors.Errorf("season with such ID found in another show")
	}

	return &Season{
		ID:           season.ID,
		SeasonNumber: season.SeasonNumber,
		ShowID:       season.ShowID,
	}, nil
}

// DeleteSeasonByID ...
func (s *Service) DeleteSeasonByID(ctx context.Context, _, seasonID uuid.UUID) error {
	if err := s.sr.DeleteSeasonByID(ctx, seasonID); err != nil {
		return fmt.Errorf("could not delete season with id=%s:%w", seasonID, err)
	}

	if err := s.sr.DeleteEpisodeBySeasonID(ctx, uuid.NullUUID{UUID: seasonID, Valid: seasonID != uuid.Nil}); err != nil {
		return fmt.Errorf("could not delete episodes with season id=%s:%w", seasonID, err)
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
	if _, err := s.sr.ReviewEpisode(ctx, repository.ReviewEpisodeParams{
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

func (s *Service) castReviewsList(ctx context.Context, source []repository.ReviewsListRow, currentUserID uuid.UUID) []Review {
	result := make([]Review, 0, len(source))
	for _, r := range source {
		review, err := s.castReview(ctx, r, currentUserID)
		if err != nil {
			log.Printf("can't cast review: %v", err)
			continue
		}
		result = append(result, review)
	}

	return result
}

func (s *Service) castReviewsListByUserID(ctx context.Context, source []repository.ReviewsListByUserIDRow, currentUserID uuid.UUID) []Review {
	result := make([]Review, 0, len(source))
	for _, r := range source {
		review, err := s.castReview(ctx, repository.ReviewsListRow(r), currentUserID)
		if err != nil {
			log.Println(err)
			continue
		}
		result = append(result, review)
	}

	return result
}

func (s *Service) castReview(ctx context.Context, source repository.ReviewsListRow, currentUserID uuid.UUID) (Review, error) {
	prof, err := s.pc.GetProfileByUserID(ctx, source.UserID, "")
	if err != nil {
		return Review{}, fmt.Errorf("could not get profile by userID=%s: %w", source.UserID, err)
	}
	username, err := s.ac.GetUsernameByID(ctx, source.UserID)
	if err != nil {
		return Review{}, fmt.Errorf("could not get username by userID=%s: %w", source.UserID, err)
	}

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
		UserAvatar: prof.Avatar,
		Rating:     int(source.Rating),
		Title:      source.Title.String,
		Review:     source.Review.String,
		CreatedAt:  source.CreatedAt.Format(time.RFC3339),
		Likes:      source.LikesNumber,
		Dislikes:   source.DislikesNumber,
		IsLiked:    isLiked,
		IsDisliked: isDisliked,
	}, nil
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

	return s.castReviewsListByUserID(ctx, reviews, currentUserID), nil
}

// LikeDislikeEpisodeReview used to store users review episode assessment (like/dislike).
func (s *Service) LikeDislikeEpisodeReview(ctx context.Context, reviewID, uid uuid.UUID, ratingType ReviewRatingType) error {
	userReview, err := s.sr.GetUserEpisodeReview(ctx, repository.GetUserEpisodeReviewParams{
		UserID:   uid,
		ReviewID: reviewID,
	})
	if err != nil && !db.IsNotFoundError(err) {
		return fmt.Errorf("could not rate episode review: %w", err)
	}

	if ReviewRatingType(userReview.RatingType.Int32) == ratingType {
		if err := s.sr.DeleteUserEpisodeReview(ctx, repository.DeleteUserEpisodeReviewParams{
			UserID:   uid,
			ReviewID: reviewID,
		}); err != nil && !db.IsNotFoundError(err) {
			return fmt.Errorf("could not unrate episode review: %w", err)
		}

		return nil
	}

	if err := s.sr.LikeDislikeEpisodeReview(ctx, repository.LikeDislikeEpisodeReviewParams{
		ReviewID: reviewID,
		UserID:   uid,
		RatingType: sql.NullInt32{
			Int32: int32(ratingType),
			Valid: true,
		}}); err != nil {
		return fmt.Errorf("could not rate episode review: %w", err)
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

// SendTipsToReviewAuthor used to send tips to an episode review author.
func (s *Service) SendTipsToReviewAuthor(ctx context.Context, reviewID, uid uuid.UUID, amount float64) error {
	review, err := s.sr.GetReviewByID(ctx, reviewID)
	if err != nil {
		return fmt.Errorf("could not get review by id: %s, error: %w", reviewID, err)
	}

	cfg := lib_solana.SendAssetsConfig{
		PercentToCharge:           s.tipsPercent,
		ChargeSolanaFeeFromSender: true,
	}
	err = s.sentTipsFunc(ctx, uid, review.UserID, amount, &cfg, fmt.Sprintf("tips for episode review: %s", reviewID))
	if err != nil {
		return fmt.Errorf("sending tips for episode review: %v, error: %w", reviewID, err)
	}

	return nil
}

// AddShowCategory ...
func (s *Service) AddShowCategory(ctx context.Context, sc ShowCategory) (ShowCategory, error) {
	category, err := s.sr.AddShowCategory(ctx, repository.AddShowCategoryParams{
		Title: sc.Title,
		Disabled: sql.NullBool{
			Bool:  sc.Disabled,
			Valid: true,
		},
		Sort: sc.Sort,
	})
	if err != nil {
		return ShowCategory{}, fmt.Errorf("could not add episode with title=%s: %w", sc.Title, err)
	}

	return castToShowCategory(category), nil
}

// DeleteShowCategoryByID ...
func (s *Service) DeleteShowCategoryByID(ctx context.Context, showCategoryID uuid.UUID) error {
	if err := s.sr.DeleteShowCategoryByID(ctx, showCategoryID); err != nil {
		return fmt.Errorf("could not delete show category with id=%s:%w", showCategoryID, err)
	}

	return nil
}

// UpdateShowCategory ...
func (s *Service) UpdateShowCategory(ctx context.Context, sc ShowCategory) error {
	if err := s.sr.UpdateShowCategory(ctx, repository.UpdateShowCategoryParams{
		ID:    sc.ID,
		Title: sc.Title,
		Disabled: sql.NullBool{
			Bool:  sc.Disabled,
			Valid: true,
		},
		Sort: sc.Sort,
	}); err != nil {
		return fmt.Errorf("could not update show category with id=%s:%w", sc.ID, err)
	}
	return nil
}

// GetShowCategoryByID returns show category with provided id.
func (s *Service) GetShowCategoryByID(ctx context.Context, showCategoryID uuid.UUID) (ShowCategory, error) {
	category, err := s.sr.GetShowCategoryByID(ctx, showCategoryID)
	if err != nil {
		return ShowCategory{}, fmt.Errorf("could not get show category with id=%s: %w", showCategoryID, err)
	}

	return castToShowCategory(category), nil
}

// GetShowCategories returns show category list.
func (s *Service) GetShowCategories(ctx context.Context, limit, offset int32) ([]ShowCategory, error) {
	category, err := s.sr.GetShowCategories(ctx, repository.GetShowCategoriesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []ShowCategory{}, fmt.Errorf("could not get show category list: %w", err)
	}

	return castToShowCategoriesList(category), nil
}

// GetShowCategoriesWithDisabled returns show category list with disabled ones.
func (s *Service) GetShowCategoriesWithDisabled(ctx context.Context, limit, offset int32) ([]ShowCategory, error) {
	category, err := s.sr.GetShowCategoriesWithDisabled(ctx, repository.GetShowCategoriesWithDisabledParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []ShowCategory{}, fmt.Errorf("could not get show category list: %w", err)
	}

	return castToShowCategoriesList(category), nil
}

// Cast []repository.ShowCategory to service []ShowCategory.
func castToShowCategoriesList(source []repository.ShowCategory) []ShowCategory {
	result := make([]ShowCategory, 0, len(source))
	for _, r := range source {
		result = append(result, castToShowCategory(r))
	}

	return result
}

// Cast repository.ShowCategory to service ShowCategory structure.
func castToShowCategory(source repository.ShowCategory) ShowCategory {
	return ShowCategory{
		ID:       source.ID,
		Title:    source.Title,
		Disabled: source.Disabled.Bool,
		Sort:     source.Sort,
	}
}

type FullEpisodeData struct {
	Episode Episode `json:"episode"`
	Season  Season  `json:"season"`
	Show    Show    `json:"show"`
}

func (s *Service) GetEpisodeByIDWithShowAndSeason(ctx context.Context, episodeID, userID uuid.UUID) (FullEpisodeData, error) {
	episode, err := s.GetPublishedEpisodeByID(ctx, episodeID, userID)
	if err != nil {
		return FullEpisodeData{}, err
	}

	season, err := s.GetSeasonByID(ctx, episode.ShowID, episode.SeasonID)
	if err != nil {
		return FullEpisodeData{}, err
	}

	show, err := s.GetPublishedShowByID(ctx, episode.ShowID)
	if err != nil {
		return FullEpisodeData{}, err
	}

	return FullEpisodeData{
		Episode: episode,
		Season:  *season,
		Show:    show,
	}, nil
}
