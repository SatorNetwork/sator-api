package puzzle_game

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/sator-api/svc/files"
	"github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

// Predefined puzzle game states
const (
	PuzzleGameStatusNew = iota
	PuzzleGameStatusInProgress
	PuzzleGameStatusFinished
)

// Predefined puzzle game results
const (
	PuzzleGameResultUndefined = iota
	PuzzleGameResultWon
	PuzzleGameResultLost
)

type (
	Service struct {
		pgr                        puzzleGameRepository
		filesSvc                   filesService
		chargeForUnlock            chargeForUnlockFunc          // function to charge user for unlocking puzzle game
		rewardsFn                  rewardsFunc                  // function to send rewards for puzzle game
		getUserRewardsMultiplierFn getUserRewardsMultiplierFunc // function to get user rewards multiplier
	}

	puzzleGameRepository interface {
		GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (repository.PuzzleGame, error)
		GetPuzzleGameByEpisodeID(ctx context.Context, episodeID uuid.UUID) (repository.PuzzleGame, error)
		CreatePuzzleGame(ctx context.Context, arg repository.CreatePuzzleGameParams) (repository.PuzzleGame, error)
		UpdatePuzzleGame(ctx context.Context, arg repository.UpdatePuzzleGameParams) (repository.PuzzleGame, error)

		GetPuzzleGameImageIDs(ctx context.Context, puzzleGameID uuid.UUID) ([]uuid.UUID, error)
		LinkImageToPuzzleGame(ctx context.Context, arg repository.LinkImageToPuzzleGameParams) error
		UnlinkImageFromPuzzleGame(ctx context.Context, arg repository.UnlinkImageFromPuzzleGameParams) error

		GetPuzzleGameCurrentAttemt(ctx context.Context, arg repository.GetPuzzleGameCurrentAttemtParams) (repository.PuzzleGamesAttempt, error)
		GetUserAvailableSteps(ctx context.Context, arg repository.GetUserAvailableStepsParams) (int32, error)
		UnlockPuzzleGame(ctx context.Context, arg repository.UnlockPuzzleGameParams) (repository.PuzzleGamesAttempt, error)
		StartPuzzleGame(ctx context.Context, arg repository.StartPuzzleGameParams) (repository.PuzzleGamesAttempt, error)
		FinishPuzzleGame(ctx context.Context, arg repository.FinishPuzzleGameParams) error

		GetPuzzleGameUnlockOption(ctx context.Context, id string) (repository.PuzzleGameUnlockOption, error)
		GetPuzzleGameUnlockOptions(ctx context.Context) ([]repository.PuzzleGameUnlockOption, error)
	}

	filesService interface {
		DeleteImageByID(ctx context.Context, id uuid.UUID) error
		GetImagesListByIDs(ctx context.Context, ids []uuid.UUID) ([]files.File, error)
	}

	chargeForUnlockFunc          func(ctx context.Context, uid uuid.UUID, amount float64, info string) error
	rewardsFunc                  func(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
	getUserRewardsMultiplierFunc func(ctx context.Context, userID uuid.UUID) (int32, error)

	// ServiceOption function
	// interface to extend service via options
	ServiceOption func(*Service)

	PuzzleGame struct {
		// general info
		ID           uuid.UUID `json:"id"`
		EpisodeID    uuid.UUID `json:"episode_id"`
		PrizePool    float64   `json:"prize_pool"`
		Rewards      float64   `json:"rewards,omitempty"`
		BonusRewards float64   `json:"bonus_rewards,omitempty"`
		PartsX       int32     `json:"parts_x"`
		// PartsY     int32     `json:"parts_y"`
		Steps      int32 `json:"steps"`
		StepsTaken int32 `json:"steps_taken,omitempty"`
		Status     int32 `json:"status"`
		Result     int32 `json:"result,omitempty"`

		// depends on user role
		Images []PuzzleGameImage `json:"images,omitempty"`
		Image  string            `json:"image,omitempty"`
	}

	PuzzleGameImage struct {
		ID      uuid.UUID `json:"id"`
		FileURL string    `json:"file_url"`
	}

	PuzzleGameUnlockOption struct {
		ID       string  `json:"id"`
		Amount   float64 `json:"amount"`
		Steps    int32   `json:"steps"`
		IsLocked bool    `json:"is_locked"`
	}
)

func NewService(pgr puzzleGameRepository, opt ...ServiceOption) *Service {
	s := &Service{pgr: pgr}

	for _, o := range opt {
		o(s)
	}

	return s
}

func NewPuzzleGameFromSQLC(pg repository.PuzzleGame) PuzzleGame {
	return PuzzleGame{
		ID:        pg.ID,
		EpisodeID: pg.EpisodeID,
		PrizePool: pg.PrizePool,
		PartsX:    pg.PartsX,
		// PartsY:    pg.PartsY,
	}
}

func (s *Service) GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (PuzzleGame, error) {
	puzzleGame, err := s.pgr.GetPuzzleGameByID(ctx, id)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game by id")
	}

	puzzleGameImages, err := s.GetPuzzleGameImages(ctx, id)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game images")
	}

	result := NewPuzzleGameFromSQLC(puzzleGame)
	result.Images = puzzleGameImages

	return result, nil
}

func (s *Service) GetPuzzleGameByEpisodeID(ctx context.Context, epID uuid.UUID) (PuzzleGame, error) {
	puzzleGame, err := s.pgr.GetPuzzleGameByEpisodeID(ctx, epID)
	if err != nil {
		puzzleGame, err = s.pgr.CreatePuzzleGame(ctx, repository.CreatePuzzleGameParams{
			EpisodeID: epID,
			PrizePool: 0,
			PartsX:    5,
		})
		if err != nil {
			return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game by episode id")
		}
	}

	puzzleGameImages, err := s.GetPuzzleGameImages(ctx, puzzleGame.ID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game images")
	}

	result := NewPuzzleGameFromSQLC(puzzleGame)
	result.Images = puzzleGameImages

	return result, nil
}

func (s *Service) CreatePuzzleGame(ctx context.Context, epID uuid.UUID, prizePool float64, partsX int32) (PuzzleGame, error) {
	puzzleGame, err := s.pgr.CreatePuzzleGame(ctx, repository.CreatePuzzleGameParams{
		EpisodeID: epID,
		PrizePool: prizePool,
		PartsX:    partsX,
	})
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't create puzzle game")
	}

	return NewPuzzleGameFromSQLC(puzzleGame), nil
}

// UpdatePuzzleGame updates puzzle game settings
func (s *Service) UpdatePuzzleGame(ctx context.Context, id uuid.UUID, prizePool float64, partsX int32) (PuzzleGame, error) {
	puzzleGame, err := s.pgr.UpdatePuzzleGame(ctx, repository.UpdatePuzzleGameParams{
		ID:        id,
		PrizePool: prizePool,
		PartsX:    partsX,
	})
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't update puzzle game")
	}

	puzzleGameImages, err := s.GetPuzzleGameImages(ctx, id)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game images")
	}

	result := NewPuzzleGameFromSQLC(puzzleGame)
	result.Images = puzzleGameImages

	return result, nil
}

// LinkImageToPuzzleGame links image to puzzle game
func (s *Service) LinkImageToPuzzleGame(ctx context.Context, gameID, fileID uuid.UUID) error {
	if err := s.pgr.LinkImageToPuzzleGame(ctx, repository.LinkImageToPuzzleGameParams{
		PuzzleGameID: gameID,
		FileID:       fileID,
	}); err != nil {
		return errors.Wrap(err, "can't link image to puzzle game")
	}

	return nil
}

// UnlinkImageFromPuzzleGame unlinks image from puzzle game
func (s *Service) UnlinkImageFromPuzzleGame(ctx context.Context, gameID, fileID uuid.UUID) error {
	if s.filesSvc != nil {
		if err := s.filesSvc.DeleteImageByID(ctx, fileID); err != nil {
			return errors.Wrap(err, "can't unlink image from puzzle game")
		}
	}

	if err := s.pgr.UnlinkImageFromPuzzleGame(ctx, repository.UnlinkImageFromPuzzleGameParams{
		PuzzleGameID: gameID,
		FileID:       fileID,
	}); err != nil {
		return errors.Wrap(err, "can't unlink image from puzzle game")
	}

	return nil
}

// GetImagesListByIDs returns list of puzzle game images by game id
func (s *Service) GetPuzzleGameImages(ctx context.Context, gameID uuid.UUID) ([]PuzzleGameImage, error) {
	if s.filesSvc == nil {
		return nil, errors.New("files service is not set")
	}

	ids, err := s.pgr.GetPuzzleGameImageIDs(ctx, gameID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get puzzle game image ids")
	}

	images, err := s.filesSvc.GetImagesListByIDs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "can't get puzzle game image urls")
	}

	result := make([]PuzzleGameImage, len(images))
	for i, img := range images {
		result[i] = PuzzleGameImage{
			ID:      img.ID,
			FileURL: img.FileUrl,
		}
	}

	return result, nil
}

func (s *Service) GetPuzzleGameForUser(ctx context.Context, userID, episodeID uuid.UUID) (PuzzleGame, error) {
	pg, err := s.pgr.GetPuzzleGameByEpisodeID(ctx, episodeID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game by episode id")
	}

	return s.getPuzzleGameForUser(ctx, userID, pg, PuzzleGameStatusNew)
}

func (s *Service) UnlockPuzzleGame(ctx context.Context, userID, puzzleGameID uuid.UUID, option string) (PuzzleGame, error) {
	if s.chargeForUnlock == nil {
		return PuzzleGame{}, errors.New("payment service is not set")
	}

	opt, err := s.pgr.GetPuzzleGameUnlockOption(ctx, option)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game unlock option")
	}

	if opt.Locked {
		return PuzzleGame{}, errors.New("this unlock option is not available")
	}

	if opt.Amount > 0 {
		if err := s.chargeForUnlock(ctx, userID, opt.Amount,
			fmt.Sprintf("Unlock puzzle game #%s", puzzleGameID.String())); err != nil {
			return PuzzleGame{}, errors.Wrap(err, "can't charge for unlock puzzle game")
		}
	}

	if _, err := s.pgr.UnlockPuzzleGame(ctx, repository.UnlockPuzzleGameParams{
		PuzzleGameID: puzzleGameID,
		UserID:       userID,
		Steps:        opt.Steps,
	}); err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't unlock puzzle game")
	}

	pg, err := s.pgr.GetPuzzleGameByID(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	result, err := s.getPuzzleGameForUser(ctx, userID, pg, PuzzleGameStatusNew)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game for user")
	}

	return result, nil
}

func (s *Service) StartPuzzleGame(ctx context.Context, userID, puzzleGameID uuid.UUID) (PuzzleGame, error) {
	img, err := s.GetRandomImageURL(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "could not start puzzle game")
	}

	if _, err := s.pgr.StartPuzzleGame(ctx, repository.StartPuzzleGameParams{
		PuzzleGameID: puzzleGameID,
		UserID:       userID,
		Image:        sql.NullString{Valid: true, String: img},
	}); err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't start puzzle game")
	}

	pg, err := s.pgr.GetPuzzleGameByID(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	result, err := s.getPuzzleGameForUser(ctx, userID, pg, PuzzleGameStatusInProgress)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game for user")
	}

	return result, nil
}

func (s *Service) FinishPuzzleGame(ctx context.Context, userID, puzzleGameID uuid.UUID, result, stepsTaken int32) (PuzzleGame, error) {
	pg, err := s.pgr.GetPuzzleGameByID(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	var rewardsAmount, lockRewardsAmount float64 = 0, 0
	if result == PuzzleGameResultWon && pg.PrizePool > 0 {
		rewardsAmount = pg.PrizePool

		if s.getUserRewardsMultiplierFn != nil {
			mltp, _ := s.getUserRewardsMultiplierFn(ctx, userID)
			if mltp > 0 {
				lockRewardsAmount = (float64(mltp) / 100) * rewardsAmount
				rewardsAmount = rewardsAmount + lockRewardsAmount
			}
		}

		if s.rewardsFn != nil {
			for i := 0; i < 5; i++ {
				err := s.rewardsFn(ctx, userID, puzzleGameID, "puzzle_games", rewardsAmount)
				if err != nil {
					log.Printf("can't add rewards for puzzle game: %v", err)
					time.Sleep(time.Second * 3)
				} else {
					break
				}
			}
		}
	}

	if err := s.pgr.FinishPuzzleGame(ctx, repository.FinishPuzzleGameParams{
		PuzzleGameID:  puzzleGameID,
		UserID:        userID,
		StepsTaken:    stepsTaken,
		RewardsAmount: rewardsAmount,
		BonusAmount:   lockRewardsAmount,
		Result:        result,
	}); err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't finish puzzle game")
	}

	pg, err = s.pgr.GetPuzzleGameByID(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	return s.getPuzzleGameForUser(ctx, userID, pg, PuzzleGameStatusFinished)
}

func (s *Service) getPuzzleGameForUser(ctx context.Context, userID uuid.UUID, puzzleGame repository.PuzzleGame, status int32) (PuzzleGame, error) {
	pg := NewPuzzleGameFromSQLC(puzzleGame)

	att, _ := s.pgr.GetPuzzleGameCurrentAttemt(ctx, repository.GetPuzzleGameCurrentAttemtParams{
		PuzzleGameID: pg.ID,
		UserID:       userID,
		Status:       status,
	})

	pg.Steps = att.Steps
	pg.StepsTaken = att.StepsTaken
	pg.Rewards = att.RewardsAmount
	pg.BonusRewards = att.BonusAmount
	pg.Result = att.Result
	pg.Status = att.Status

	if !att.Image.Valid {
		img, err := s.GetRandomImageURL(ctx, pg.ID)
		if err != nil {
			return PuzzleGame{}, err
		}
		pg.Image = img
	} else {
		pg.Image = att.Image.String
	}

	return pg, nil
}

func (s *Service) GetRandomImageURL(ctx context.Context, gameID uuid.UUID) (string, error) {
	images, err := s.GetPuzzleGameImages(ctx, gameID)
	if err != nil {
		return "", errors.Wrap(err, "can't get puzzle game images")
	}

	img := getRandomImage(images)
	if img == "" {
		return "", fmt.Errorf("no images for puzzle game %s", gameID)
	}

	return img, nil
}

func getRandomImage(images []PuzzleGameImage) string {
	if len(images) == 0 {
		return ""
	}

	rand.Seed(time.Now().Unix())
	imgKey := rand.Intn(len(images))
	if img := images[imgKey]; img.ID != uuid.Nil {
		return img.FileURL
	}

	return getRandomImage(images)
}

// GetPuzzleGameUnlockOptions returns all available puzzle game unlock options
func (s *Service) GetPuzzleGameUnlockOptions(ctx context.Context) ([]PuzzleGameUnlockOption, error) {
	options, err := s.pgr.GetPuzzleGameUnlockOptions(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "can't get puzzle game unlock options")
	}

	result := make([]PuzzleGameUnlockOption, 0, len(options))
	for _, opt := range options {
		result = append(result, PuzzleGameUnlockOption{
			ID:       opt.ID,
			Amount:   opt.Amount,
			Steps:    opt.Steps,
			IsLocked: opt.Locked,
		})
	}

	return result, nil
}
