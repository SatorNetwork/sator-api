package puzzle_game

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/SatorNetwork/gopuzzlegame"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/svc/files"
	"github.com/SatorNetwork/sator-api/svc/puzzle_game/repository"
)

// Predefined puzzle game states
const (
	PuzzleGameStatusNew = iota
	PuzzleGameStatusInProgress
	PuzzleGameStatusFinished
	PuzzleStatusReachedStepLimit
)

type (
	Service struct {
		pgr                        puzzleGameRepository
		filesSvc                   filesService
		chargeForUnlock            chargeForUnlockFunc          // function to charge user for unlocking puzzle game
		rewardsFn                  rewardsFunc                  // function to send rewards for puzzle game
		getUserRewardsMultiplierFn getUserRewardsMultiplierFunc // function to get user rewards multiplier
		puzzleGameShuffle          bool
	}

	puzzleGameRepository interface {
		GetPuzzleGameByID(ctx context.Context, id uuid.UUID) (repository.PuzzleGame, error)
		GetPuzzleGameByEpisodeID(ctx context.Context, episodeID uuid.UUID) (repository.PuzzleGame, error)
		CreatePuzzleGame(ctx context.Context, arg repository.CreatePuzzleGameParams) (repository.PuzzleGame, error)
		UpdatePuzzleGame(ctx context.Context, arg repository.UpdatePuzzleGameParams) (repository.PuzzleGame, error)

		GetPuzzleGameImageIDs(ctx context.Context, puzzleGameID uuid.UUID) ([]uuid.UUID, error)
		LinkImageToPuzzleGame(ctx context.Context, arg repository.LinkImageToPuzzleGameParams) error
		UnlinkImageFromPuzzleGame(ctx context.Context, arg repository.UnlinkImageFromPuzzleGameParams) error

		GetPuzzleGameCurrentAttempt(ctx context.Context, arg repository.GetPuzzleGameCurrentAttemptParams) (repository.PuzzleGamesAttempt, error)
		UpdatePuzzleGameAttempt(ctx context.Context, arg repository.UpdatePuzzleGameAttemptParams) (repository.PuzzleGamesAttempt, error)
		GetUserAvailableSteps(ctx context.Context, arg repository.GetUserAvailableStepsParams) (int32, error)
		FinishPuzzleGame(ctx context.Context, arg repository.FinishPuzzleGameParams) error
		UnlockPuzzleGame(ctx context.Context, arg repository.UnlockPuzzleGameParams) (repository.PuzzleGamesAttempt, error)
		StartPuzzleGame(ctx context.Context, arg repository.StartPuzzleGameParams) (repository.PuzzleGamesAttempt, error)

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
		Steps      int32                `json:"steps"`
		StepsTaken int32                `json:"steps_taken,omitempty"`
		Status     int32                `json:"status"`
		Tiles      []*gopuzzlegame.Tile `json:"tiles,omitempty"`

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

func (pg *PuzzleGame) HideCorrectPositions() PuzzleGame {
	for _, tile := range pg.Tiles {
		tile.CorrectPosition = nil
	}
	return *pg
}

func NewService(pgr puzzleGameRepository, puzzleGameShuffle bool, opt ...ServiceOption) *Service {
	s := &Service{pgr: pgr, puzzleGameShuffle: puzzleGameShuffle}

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

func (s *Service) GetPuzzleGameByEpisodeID(ctx context.Context, epID uuid.UUID, isTestUser bool) (PuzzleGame, error) {
	if isTestUser {
		return PuzzleGame{}, ErrNotFound
	}

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

func (s *Service) GetPuzzleGameForUser(ctx context.Context, userID, episodeID uuid.UUID, isTestUser bool) (PuzzleGame, error) {
	if isTestUser {
		return PuzzleGame{}, ErrNotFound
	}

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

	pg, err := s.pgr.GetPuzzleGameByID(ctx, puzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	p, err := gopuzzlegame.GeneratePuzzle(int(pg.PartsX), s.puzzleGameShuffle)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't generate puzzle")
	}

	rawTiles, err := json.Marshal(p.Tiles)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't create puzzle game")
	}

	if _, err := s.pgr.StartPuzzleGame(ctx, repository.StartPuzzleGameParams{
		PuzzleGameID: puzzleGameID,
		UserID:       userID,
		Image:        sql.NullString{Valid: true, String: img},
		Tiles:        sql.NullString{String: string(rawTiles), Valid: true},
	}); err != nil {
		fmt.Println(err.Error())
		return PuzzleGame{}, errors.Wrap(err, "can't start puzzle game")
	}

	result, err := s.getPuzzleGameForUser(ctx, userID, pg, PuzzleGameStatusInProgress)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game for user")
	}

	result.Tiles = p.Tiles
	return result.HideCorrectPositions(), nil
}

func (s *Service) getPuzzleGameForUser(ctx context.Context, userID uuid.UUID, puzzleGame repository.PuzzleGame, status int32) (PuzzleGame, error) {
	pg := NewPuzzleGameFromSQLC(puzzleGame)

	att, _ := s.pgr.GetPuzzleGameCurrentAttempt(ctx, repository.GetPuzzleGameCurrentAttemptParams{
		PuzzleGameID: pg.ID,
		UserID:       userID,
		Status:       status,
	})

	pg.Steps = att.Steps
	pg.StepsTaken = att.StepsTaken
	pg.Rewards = att.RewardsAmount
	pg.BonusRewards = att.BonusAmount
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

func (s *Service) TapTile(ctx context.Context, userID, puzzleGameID uuid.UUID, position gopuzzlegame.Position) (PuzzleGame, error) {
	att, err := s.pgr.GetPuzzleGameCurrentAttempt(ctx, repository.GetPuzzleGameCurrentAttemptParams{
		UserID:       userID,
		PuzzleGameID: puzzleGameID,
		Status:       PuzzleGameStatusInProgress,
	})
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game attempt")
	}

	if att.Status != PuzzleGameStatusInProgress || att.Steps == 0 {
		return PuzzleGame{}, errors.New("puzzle game is not in progress")
	}

	if att.StepsTaken >= att.Steps || att.RewardsAmount > 0 {
		return PuzzleGame{}, errors.New("puzzle game is over")
	}

	pg, err := s.pgr.GetPuzzleGameByID(ctx, att.PuzzleGameID)
	if err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game")
	}

	var tiles []*gopuzzlegame.Tile
	if err = json.Unmarshal([]byte(att.Tiles.String), &tiles); err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't get tiles")
	}

	controller := gopuzzlegame.PuzzleController{
		PuzzleStatus: att.Status,
		Puzzle:       &gopuzzlegame.Puzzle{Tiles: tiles},
		StepsTaken:   att.StepsTaken,
		Steps:        att.Steps,
	}

	var tile *gopuzzlegame.Tile
	for _, t := range tiles {
		if t.CurrentPosition == position {
			tile = t
			break
		}
	}
	if tile == nil {
		return PuzzleGame{}, errors.New("can't get tile with such position")
	}

	if err = controller.TapTile(tile); err != nil {
		return PuzzleGame{}, errors.Wrap(err, "can't tap tile")
	}

	tilesBytes, err := json.Marshal(controller.Puzzle.Tiles)
	if err != nil {
		return PuzzleGame{}, err
	}

	att.StepsTaken = controller.StepsTaken
	att.Status = controller.PuzzleStatus
	att.Tiles = sql.NullString{String: string(tilesBytes), Valid: true}

	var rewardsAmount, lockRewardsAmount float64 = 0, 0
	if att.Status == PuzzleGameStatusFinished {
		if pg.PrizePool > 0 {
			rewardsAmount = pg.PrizePool

			if s.getUserRewardsMultiplierFn != nil {
				mltp, _ := s.getUserRewardsMultiplierFn(ctx, userID)
				if mltp > 0 {
					lockRewardsAmount = (float64(mltp) / 100) * rewardsAmount
					rewardsAmount = rewardsAmount + lockRewardsAmount
				}
			}

			if s.rewardsFn != nil {
				att, err := s.pgr.GetPuzzleGameCurrentAttempt(ctx, repository.GetPuzzleGameCurrentAttemptParams{
					PuzzleGameID: puzzleGameID,
					UserID:       userID,
					Status:       PuzzleGameStatusInProgress,
				})
				if err != nil {
					return PuzzleGame{}, errors.Wrap(err, "can't get puzzle game current attempt")
				}

			retrySendReward:
				for i := 0; i < 5; i++ {
					if err := s.rewardsFn(ctx, userID, att.ID, "puzzle_games", rewardsAmount); err != nil {
						log.Printf("can't add rewards for puzzle game: %v", err)
						time.Sleep(time.Second * 3)
					} else {
						break retrySendReward
					}
				}
			}
		}

		if err := s.pgr.FinishPuzzleGame(ctx, repository.FinishPuzzleGameParams{
			PuzzleGameID:  puzzleGameID,
			UserID:        userID,
			StepsTaken:    att.StepsTaken,
			RewardsAmount: rewardsAmount,
			BonusAmount:   lockRewardsAmount,
		}); err != nil {
			return PuzzleGame{}, errors.Wrap(err, "can't finish puzzle game")
		}
	} else {
		_, err = s.pgr.UpdatePuzzleGameAttempt(ctx, repository.UpdatePuzzleGameAttemptParams{
			Status:     att.Status,
			Steps:      att.Steps,
			StepsTaken: att.StepsTaken,
			Tiles:      sql.NullString{String: string(tilesBytes), Valid: true},
			ID:         att.ID,
		})
		if err != nil {
			return PuzzleGame{}, errors.Wrap(err, "can't update game attempt")
		}
	}

	response := PuzzleGame{
		ID:           pg.ID,
		EpisodeID:    pg.EpisodeID,
		PrizePool:    pg.PrizePool,
		Rewards:      rewardsAmount,
		BonusRewards: lockRewardsAmount,
		PartsX:       pg.PartsX,
		Steps:        att.Steps,
		StepsTaken:   att.StepsTaken,
		Status:       att.Status,
		Tiles:        controller.Puzzle.Tiles,
		Images:       nil,
		Image:        att.Image.String,
	}
	return response.HideCorrectPositions(), nil
}
