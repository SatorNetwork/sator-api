package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/SatorNetwork/sator-api/lib/db"
	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		gameRepo             gameRepository
		minVersion           string
		energyFull           int32
		energyRecoveryPeriod time.Duration
		minRewardsToClaim    float64
		rewards              rewardsService
	}

	ServiceOption func(*Service)

	gameRepository interface {
		StartGame(ctx context.Context, arg repository.StartGameParams) error
		GetCurrentGame(ctx context.Context, userID uuid.UUID) (repository.UnityGameResult, error)
		FinishGame(ctx context.Context, arg repository.FinishGameParams) error

		AddNFTPack(ctx context.Context, arg repository.AddNFTPackParams) error
		GetNFTPack(ctx context.Context, id uuid.UUID) (repository.UnityGameNftPack, error)
		GetNFTPacksList(ctx context.Context) ([]repository.UnityGameNftPack, error)
		UpdateNFTPack(ctx context.Context, arg repository.UpdateNFTPackParams) error
		SoftDeleteNFTPack(ctx context.Context, id uuid.UUID) error
		DeleteNFTPack(ctx context.Context, id uuid.UUID) error

		AddNFT(ctx context.Context, arg repository.AddNFTParams) (repository.UnityGameNft, error)
		GetNFT(ctx context.Context, id string) (repository.UnityGameNft, error)
		GetNFTs(ctx context.Context, arg repository.GetNFTsParams) ([]repository.UnityGameNft, error)
		GetNFTsByTypeAndLevel(ctx context.Context, arg repository.GetNFTsByTypeAndLevelParams) ([]repository.UnityGameNft, error)
		UpdateNFT(ctx context.Context, arg repository.UpdateNFTParams) error
		SoftDeleteNFT(ctx context.Context, id string) error
		DeleteNFT(ctx context.Context, id string) error

		CraftNFTs(ctx context.Context, arg repository.CraftNFTsParams) error
		GetNFTsByPlayer(ctx context.Context, userID uuid.UUID) ([]repository.UnityGameNft, error)
		LinkNFTToPlayer(ctx context.Context, arg repository.LinkNFTToPlayerParams) error
		UnlinkNFTFromPlayer(ctx context.Context, arg repository.UnlinkNFTFromPlayerParams) error

		AddNewPlayer(ctx context.Context, arg repository.AddNewPlayerParams) (repository.UnityGamePlayer, error)
		GetPlayer(ctx context.Context, userID uuid.UUID) (repository.UnityGamePlayer, error)
		RefillEnergyOfPlayer(ctx context.Context, arg repository.RefillEnergyOfPlayerParams) error
		StoreSelectedNFT(ctx context.Context, arg repository.StoreSelectedNFTParams) error
	}

	rewardsService interface {
		AddDepositTransaction(ctx context.Context, userID, relationID uuid.UUID, relationType string, amount float64) error
		GetUserRewards(ctx context.Context, userID uuid.UUID) (total float64, available float64, err error)
	}

	PlayerInfo struct {
		UserID        uuid.UUID
		EnergyPoints  int
		SelectedNftID string
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo gameRepository, rs rewardsService, opt ...ServiceOption) *Service {
	s := &Service{
		gameRepo:             repo,
		energyFull:           3,
		energyRecoveryPeriod: time.Hour * 4,
		minRewardsToClaim:    50,
		rewards:              rs,
	}

	// Apply options
	for _, o := range opt {
		o(s)
	}

	return s
}

// GetPlayerInfo ...
func (s *Service) GetPlayerInfo(ctx context.Context, uid uuid.UUID) (*PlayerInfo, error) {
	player, err := s.gameRepo.GetPlayer(ctx, uid)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, err
		}

		player, err = s.gameRepo.AddNewPlayer(ctx, repository.AddNewPlayerParams{
			UserID:           uid,
			EnergyPoints:     s.energyFull,
			EnergyRefilledAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	energy := player.EnergyPoints
	if player.EnergyPoints < s.energyFull && time.Since(player.EnergyRefilledAt) > s.energyRecoveryPeriod {
		hoursSince := time.Since(player.EnergyRefilledAt).Hours()
		recoveryHours := s.energyRecoveryPeriod.Hours()
		recoveryPoints := math.Floor(hoursSince / recoveryHours)
		if recoveryPoints > 0 {
			energy = player.EnergyPoints + int32(recoveryPoints)
			if energy > s.energyFull {
				energy = s.energyFull
			}
			if err := s.gameRepo.RefillEnergyOfPlayer(ctx, repository.RefillEnergyOfPlayerParams{
				UserID:           uid,
				EnergyPoints:     energy,
				EnergyRefilledAt: time.Now(),
			}); err != nil {
				return nil, err
			}
		}
	}

	return &PlayerInfo{
		UserID:        player.UserID,
		EnergyPoints:  int(energy),
		SelectedNftID: player.SelectedNftID.String,
	}, nil
}

// GetMinVersion ...
func (s *Service) GetMinVersion(ctx context.Context) string {
	return s.minVersion
}

// GetEnergyLeft ...
func (s *Service) GetEnergyLeft(ctx context.Context, uid uuid.UUID) (int, error) {
	player, err := s.GetPlayerInfo(ctx, uid)
	if err != nil {
		return 0, err
	}

	return player.EnergyPoints, nil
}

// GetUserNFTs ...
func (s *Service) GetUserNFTs(ctx context.Context, uid uuid.UUID) ([]NFTInfo, error) {
	nfts, err := s.gameRepo.GetNFTsByPlayer(ctx, uid)
	if err != nil {
		return nil, err
	}

	result := make([]NFTInfo, 0, len(nfts))
	for _, nft := range nfts {
		result = append(result, NFTInfo{
			ID:       nft.ID,
			MaxLevel: nft.AllowedLevels,
			NftType:  NFTType(nft.NftType),
		})
	}

	return result, nil
}

// GetSelectedNFT ...
func (s *Service) GetSelectedNFT(ctx context.Context, uid uuid.UUID) (*string, error) {
	player, err := s.GetPlayerInfo(ctx, uid)
	if err != nil {
		return nil, err
	}

	if player.SelectedNftID == "" {
		return nil, nil
	}

	return &player.SelectedNftID, nil
}

// GetNFTPacks ...
func (s *Service) GetNFTPacks(ctx context.Context, uid uuid.UUID) ([]NFTPackInfo, error) {
	nftPacks, err := s.gameRepo.GetNFTPacksList(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]NFTPackInfo, 0, len(nftPacks))
	for _, nftPack := range nftPacks {
		dropChances := DropChances{}
		if err := json.Unmarshal(nftPack.DropChances, &dropChances); err != nil {
			return nil, err
		}

		result = append(result, NFTPackInfo{
			ID:          nftPack.ID.String(),
			DropChances: dropChances,
			Price:       nftPack.Price,
		})
	}

	return result, nil
}

// BuyNFTPack ...
func (s *Service) BuyNFTPack(ctx context.Context, uid, packID uuid.UUID) error {
	return nil
}

// CraftNFT ...
func (s *Service) CraftNFT(ctx context.Context, uid uuid.UUID, nftsToCraft []string) error {
	return nil
}

// SelectNFT ...
func (s *Service) SelectNFT(ctx context.Context, uid uuid.UUID, nftMintAddr string) error {
	if err := s.gameRepo.StoreSelectedNFT(ctx, repository.StoreSelectedNFTParams{
		UserID:        uid,
		SelectedNftID: sql.NullString{String: nftMintAddr, Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to store selected nft: %w", err)
	}
	return nil
}

// StartGame ...
func (s *Service) StartGame(ctx context.Context, uid uuid.UUID, complexity string, isTraining bool) error {
	player, err := s.GetPlayerInfo(ctx, uid)
	if err != nil {
		return err
	}

	if err := s.gameRepo.StartGame(ctx, repository.StartGameParams{
		UserID:     uid,
		NftID:      player.SelectedNftID,
		Complexity: complexity,
		IsTraining: isTraining,
	}); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	return nil
}

// FinishGame ...
// TODO: wrapp to db transaction
// TODO: rewards calculation
func (s *Service) FinishGame(ctx context.Context, uid uuid.UUID, blocksDone int32) error {
	currentGame, err := s.gameRepo.GetCurrentGame(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get current game: %w", err)
	}

	var rewardsAmount float64 = float64(blocksDone) * 0.5

	if err := s.gameRepo.FinishGame(ctx, repository.FinishGameParams{
		ID:         currentGame.ID,
		BlocksDone: blocksDone,
		Rewards:    rewardsAmount,
	}); err != nil {
		return fmt.Errorf("failed to finish game: %w", err)
	}

	if err := s.rewards.AddDepositTransaction(ctx, uid, currentGame.ID, "game", rewardsAmount); err != nil {
		return fmt.Errorf("failed to add rewards: %w", err)
	}

	return nil
}

// GetDefaultGameConfig ...
func (s *Service) GetDefaultGameConfig(ctx context.Context, uid uuid.UUID) (*GameConfig, error) {
	cnf := &GameConfig{}
	if err := json.Unmarshal([]byte(defaultGameConfig), cnf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game config: %w", err)
	}

	return cnf, nil
}

// GetMinAmountToClaim ...
func (s *Service) GetMinAmountToClaim() float64 {
	return s.minRewardsToClaim
}

// GetUserRewards ...
func (s *Service) GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error) {
	total, _, err := s.rewards.GetUserRewards(ctx, uid)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// claimRewardsFunc is a function that can be used to claim rewards to the user's SAO wallet
type claimRewardsFunc func(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error)

// ClaimRewards ...
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64, claimFn claimRewardsFunc) error {
	userRewardsAmount, _ := s.GetUserRewards(ctx, uid)
	if userRewardsAmount < s.GetMinAmountToClaim() {
		return fmt.Errorf("not enough rewards to claim, need %f, have %f", s.GetMinAmountToClaim(), userRewardsAmount)
	}
	if amount > userRewardsAmount {
		return fmt.Errorf("not enough rewards to claim, need %f, have %f", amount, userRewardsAmount)
	}

	_, err := claimFn(ctx, uid, userRewardsAmount)
	if err != nil {
		return fmt.Errorf("failed to claim rewards: %w", err)
	}

	return nil
}
