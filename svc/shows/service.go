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
		CategoryID    uuid.UUID `json:"category"`
		Description   string    `json:"description"`
	}

	// Episode struct
	Episode struct {
		ID            uuid.UUID `json:"id"`
		ShowID        uuid.UUID `json:"show_id"`
		EpisodeNumber int32     `json:"episode_number"`
		Cover         string    `json:"cover"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		ReleaseDate   string    `json:"release_date"`
	}

	// ShowCategory struct
	ShowCategory struct {
		ID           uuid.UUID `json:"id"`
		CategoryName string    `json:"category_name"`
		Title        string    `json:"title"`
		Disabled     bool      `json:"disabled"`
	}

	HomePage struct {
		CategoryName string `json:"category_name"`
		Title        string `json:"title"`
		URL          string `json:"url"`
		ListShows    []Show `json:"list_shows"`
	}

	showsRepository interface {
		AddShow(ctx context.Context, arg repository.AddShowParams) (repository.AddShowRow, error)
		DeleteShowByID(ctx context.Context, id uuid.UUID) error
		GetShows(ctx context.Context) ([]repository.GetShowsRow, error)
		GetShowsPaginated(ctx context.Context, arg repository.GetShowsPaginatedParams) ([]repository.GetShowsPaginatedRow, error)
		GetShowByID(ctx context.Context, id uuid.UUID) (repository.GetShowByIDRow, error)
		UpdateShow(ctx context.Context, arg repository.UpdateShowParams) error

		AddEpisode(ctx context.Context, arg repository.AddEpisodeParams) (repository.Episode, error)
		GetEpisodeByID(ctx context.Context, arg repository.GetEpisodeByIDParams) (repository.Episode, error)
		GetEpisodesByShowID(ctx context.Context, arg repository.GetEpisodesByShowIDParams) ([]repository.Episode, error)
		DeleteEpisodeByID(ctx context.Context, arg repository.DeleteEpisodeByIDParams) error
		UpdateEpisode(ctx context.Context, arg repository.UpdateEpisodeParams) error

		AddShowCategory(ctx context.Context, arg repository.AddShowCategoryParams) (repository.ShowsCategory, error)
		DeleteShowCategoryByID(ctx context.Context, id uuid.UUID) error
		GetShowCategoryByID(ctx context.Context, id uuid.UUID) (repository.ShowsCategory, error)
		GetShowCategories(ctx context.Context) ([]repository.ShowsCategory, error)
		UpdateShowCategory(ctx context.Context, arg repository.UpdateShowCategoryParams) error

		AddShowToCategory(ctx context.Context, arg repository.AddShowToCategoryParams) error
		DeleteShowToCategory(ctx context.Context, arg repository.DeleteShowToCategoryParams) error
		DeleteShowToCategoryByShowID(ctx context.Context, showID uuid.UUID) error
		UpdateShowToCategory(ctx context.Context, arg repository.UpdateShowToCategoryParams) error
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
func (s *Service) GetShows(ctx context.Context, limit, offset int32) ([]Show, error) {
	shows, err := s.sr.GetShowsPaginated(ctx, repository.GetShowsPaginatedParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return []Show{}, fmt.Errorf("could not get shows list: %w", err)
	}

	return castToListShow(shows), nil
}

// GetShowsHome returns shows.
func (s *Service) GetShowsHome(ctx context.Context) (homePage []HomePage, err error) {
	shows, err := s.sr.GetShows(ctx)
	if err != nil {
		return []HomePage{}, fmt.Errorf("could not get shows list: %w", err)
	}

	categories, err := s.sr.GetShowCategories(ctx)
	if err != nil {
		return []HomePage{}, fmt.Errorf("could not get categories list: %w", err)
	}

	homePage = append(homePage, HomePage{
		CategoryName: "All",
		Title:        "All",
		URL:          "/categories/all",
		ListShows:    castGetShowsRowToShow(shows),
	})

	for _, category := range categories {
		var listShows []Show
		for _, show := range shows {
			if show.CategoryID == category.ID {
				listShows = append(listShows, castGetShowRowToShow(show))
			}
		}
		homePage = append(homePage, HomePage{
			CategoryName: category.CategoryName,
			Title:        category.Title,
			URL:          "/category/" + category.Title,
			ListShows:    listShows,
		})
	}

	return homePage, nil
}

// castGetShowsRowToShow cast repository.GetShowsRow to service Show structure
func castGetShowRowToShow(s repository.GetShowsRow) Show {
	return Show{
		ID:            s.ID,
		Title:         s.Title,
		Cover:         s.Cover,
		HasNewEpisode: s.HasNewEpisode,
		Description:   s.Description.String,
		CategoryID:    s.CategoryID,
	}
}

// castGetShowsRowToShow cast repository.GetShowsRow to service []Show structure
func castGetShowsRowToShow(source []repository.GetShowsRow) []Show {
	result := make([]Show, 0, len(source))
	for _, s := range source {
		result = append(result, Show{
			ID:            s.ID,
			Title:         s.Title,
			Cover:         s.Cover,
			HasNewEpisode: s.HasNewEpisode,
			Description:   s.Description.String,
			CategoryID:    s.CategoryID,
		})
	}

	return result
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
func castToListShow(source []repository.GetShowsPaginatedRow) []Show {
	result := make([]Show, 0, len(source))
	for _, s := range source {
		result = append(result, Show{
			ID:            s.ID,
			Title:         s.Title,
			Cover:         s.Cover,
			HasNewEpisode: s.HasNewEpisode,
			Description:   s.Description.String,
			CategoryID:    s.CategoryID,
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

// Cast repository.GetShowByIDRow to service Show structure
func castToShow(s repository.GetShowByIDRow) Show {
	return Show{
		ID:            s.ID,
		Title:         s.Title,
		Cover:         s.Cover,
		HasNewEpisode: s.HasNewEpisode,
		Description:   s.Description.String,
		CategoryID:    s.CategoryID,
	}
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
	}
}

// AddShow ..
func (s *Service) AddShow(ctx context.Context, sh Show) (Show, error) {
	show, err := s.sr.AddShow(ctx, repository.AddShowParams{
		Title:         sh.Title,
		Cover:         sh.Cover,
		HasNewEpisode: sh.HasNewEpisode,
		Description: sql.NullString{
			String: sh.Description,
			Valid:  len(sh.Description) > 0,
		},
	})
	if err != nil {
		return Show{}, fmt.Errorf("could not add show with title=%s: %w", sh.Title, err)
	}

	err = s.sr.AddShowToCategory(ctx, repository.AddShowToCategoryParams{
		CategoryID: sh.CategoryID,
		ShowID:     show.ID,
	})
	if err != nil {
		return Show{}, fmt.Errorf("could not add show to category with title=%s: %w", sh.Title, err)
	}

	return castAddShowRowToShow(show), nil
}

// Cast repository.AddShowRow to service Show structure
func castAddShowRowToShow(s repository.AddShowRow) Show {
	return Show{
		ID:            s.ID,
		Title:         s.Title,
		Cover:         s.Cover,
		HasNewEpisode: s.HasNewEpisode,
		Description:   s.Description.String,
		CategoryID:    s.CategoryID,
	}
}

// UpdateShow ..
func (s *Service) UpdateShow(ctx context.Context, sh Show) error {
	err := s.sr.UpdateShow(ctx, repository.UpdateShowParams{
		Title:         sh.Title,
		Cover:         sh.Cover,
		HasNewEpisode: sh.HasNewEpisode,
		Description: sql.NullString{
			String: sh.Description,
			Valid:  len(sh.Description) > 0,
		},
		ID: sh.ID,
	})
	if err != nil {
		return fmt.Errorf("could not update show with id=%s:%w", sh.ID, err)
	}

	err = s.sr.UpdateShowToCategory(ctx, repository.UpdateShowToCategoryParams{
		CategoryID: sh.CategoryID,
		ShowID:     sh.ID,
	})
	if err != nil {
		return fmt.Errorf("could not update show to category with title=%s: %w", sh.Title, err)
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
	})
	if err != nil {
		return Episode{}, fmt.Errorf("could not add episode with show_id=%s, episodeNumber=%v: %w", ep.ShowID, ep.EpisodeNumber, err)
	}

	return castToEpisode(episode), nil
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

// AddShowCategories ..
func (s *Service) AddShowCategories(ctx context.Context, sc ShowCategory) (ShowCategory, error) {
	category, err := s.sr.AddShowCategory(ctx, repository.AddShowCategoryParams{
		CategoryName: sc.CategoryName,
		Title:        sc.Title,
		Disabled: sql.NullBool{
			Bool:  sc.Disabled,
			Valid: true,
		},
	})
	if err != nil {
		return ShowCategory{}, fmt.Errorf("could not add episode with category_name=%s: %w", sc.CategoryName, err)
	}

	return castToShowCategory(category), nil
}

// Cast repository.ShowsCategory to service ShowCategory structure
func castToShowCategory(source repository.ShowsCategory) ShowCategory {
	return ShowCategory{
		ID:           source.ID,
		CategoryName: source.CategoryName,
		Title:        source.Title,
		Disabled:     source.Disabled.Bool,
	}
}

// DeleteShowCategoryByID ..
func (s *Service) DeleteShowCategoryByID(ctx context.Context, showCategoryID uuid.UUID) error {
	if err := s.sr.DeleteShowCategoryByID(ctx, showCategoryID); err != nil {
		return fmt.Errorf("could not delete show category with id=%s:%w", showCategoryID, err)
	}

	return nil
}

// UpdateShowCategory ..
func (s *Service) UpdateShowCategory(ctx context.Context, sc ShowCategory) error {
	if err := s.sr.UpdateShowCategory(ctx, repository.UpdateShowCategoryParams{
		ID:           sc.ID,
		CategoryName: sc.CategoryName,
		Title:        sc.Title,
		Disabled: sql.NullBool{
			Bool:  sc.Disabled,
			Valid: true,
		},
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

// GetShowCategories returns categories.
func (s *Service) GetShowCategories(ctx context.Context) ([]ShowCategory, error) {
	categories, err := s.sr.GetShowCategories(ctx)
	if err != nil {
		return []ShowCategory{}, fmt.Errorf("could not get categories list: %w", err)
	}

	return castToListShowCategory(categories), nil
}

// Cast repository.ShowsCategory to service ShowCategory structure
func castToListShowCategory(source []repository.ShowsCategory) []ShowCategory {
	result := make([]ShowCategory, 0, len(source))
	for _, s := range source {
		result = append(result, ShowCategory{
			ID:           s.ID,
			CategoryName: s.CategoryName,
			Title:        s.Title,
			Disabled:     s.Disabled.Bool,
		})
	}

	return result
}

// AddShowToCategory ..
func (s *Service) AddShowToCategory(ctx context.Context, categoryID, showID uuid.UUID) error {
	if err := s.sr.AddShowToCategory(ctx, repository.AddShowToCategoryParams{
		CategoryID: categoryID,
		ShowID:     showID,
	}); err != nil {
		return fmt.Errorf("could not add show to category with category_id=%s, show_id=%s: %w", categoryID, showID, err)
	}

	return nil
}

// DeleteShowToCategory ..
func (s *Service) DeleteShowToCategory(ctx context.Context, categoryID, showID uuid.UUID) error {
	if err := s.sr.DeleteShowToCategory(ctx, repository.DeleteShowToCategoryParams{
		CategoryID: categoryID,
		ShowID:     showID,
	}); err != nil {
		return fmt.Errorf("could not delete show to category with category_id=%s, show_id=%s: %w", categoryID, showID, err)
	}

	return nil
}

// DeleteShowToCategoryByShowID ..
func (s *Service) DeleteShowToCategoryByShowID(ctx context.Context, showID uuid.UUID) error {
	if err := s.sr.DeleteShowToCategoryByShowID(ctx, showID); err != nil {
		return fmt.Errorf("could not delete show to categories with show_id=%s: %w", showID, err)
	}

	return nil
}
