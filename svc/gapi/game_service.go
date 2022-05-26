package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/SatorNetwork/sator-api/lib/db"
	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		db                   *sql.DB
		gameRepo             gameRepository
		conf                 configer
		minVersion           string
		energyFull           int32
		energyRecoveryPeriod time.Duration
		minRewardsToClaim    float64
		craftStepAmount      float64
	}

	ServiceOption func(*Service)

	gameRepository interface {
		WithTx(tx *sql.Tx) *repository.Queries

		StartGame(ctx context.Context, arg repository.StartGameParams) error
		GetCurrentGame(ctx context.Context, userID uuid.UUID) (repository.UnityGameResult, error)
		FinishGame(ctx context.Context, arg repository.FinishGameParams) error

		GetNFTPacksList(ctx context.Context) ([]repository.UnityGameNftPack, error)

		AddNFT(ctx context.Context, arg repository.AddNFTParams) (repository.UnityGameNft, error)
		GetUserNFT(ctx context.Context, arg repository.GetUserNFTParams) (repository.UnityGameNft, error)
		GetUserNFTByIDs(ctx context.Context, arg repository.GetUserNFTByIDsParams) ([]repository.UnityGameNft, error)
		CraftNFTs(ctx context.Context, arg repository.CraftNFTsParams) error
		GetUserNFTs(ctx context.Context, userID uuid.UUID) ([]repository.UnityGameNft, error)

		AddNewPlayer(ctx context.Context, arg repository.AddNewPlayerParams) (repository.UnityGamePlayer, error)
		GetPlayer(ctx context.Context, userID uuid.UUID) (repository.UnityGamePlayer, error)
		RefillEnergyOfPlayer(ctx context.Context, arg repository.RefillEnergyOfPlayerParams) error
		StoreSelectedNFT(ctx context.Context, arg repository.StoreSelectedNFTParams) error

		GetUserRewards(ctx context.Context, userID uuid.UUID) (float64, error)
		RewardsDeposit(ctx context.Context, arg repository.RewardsDepositParams) error
		RewardsWithdraw(ctx context.Context, arg repository.RewardsWithdrawParams) error
		GetUserRewardsDeposited(ctx context.Context, userID uuid.UUID) (float64, error)
		GetUserRewardsWithdrawn(ctx context.Context, userID uuid.UUID) (float64, error)
	}

	configer interface {
		GetBool(ctx context.Context, key string) (bool, error)
		GetString(ctx context.Context, key string) (string, error)
		GetFloat64(ctx context.Context, key string) (float64, error)
		GetInt(ctx context.Context, key string) (int, error)
		GetJSON(ctx context.Context, key string, result interface{}) error
		GetDurration(ctx context.Context, key string) (time.Duration, error)
	}

	PlayerInfo struct {
		UserID        uuid.UUID
		EnergyPoints  int
		SelectedNftID string
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo gameRepository, conf configer, opt ...ServiceOption) *Service {
	s := &Service{
		gameRepo:             repo,
		conf:                 conf,
		energyFull:           3,
		energyRecoveryPeriod: time.Hour * 4,
		minRewardsToClaim:    50,
		craftStepAmount:      500,
	}

	// Apply options
	for _, o := range opt {
		o(s)
	}

	if energyFull, err := conf.GetInt(context.Background(), "energy_full"); err != nil {
		log.Printf("[WARN] energy_full not found in config, using default value: %d", s.energyFull)
	} else {
		s.energyFull = int32(energyFull)
	}

	if energyRecoveryPeriod, err := conf.GetDurration(context.Background(), "energy_recovery_period"); err != nil {
		log.Printf("[WARN] energy_recovery_period not found in config, using default value: %s", s.energyRecoveryPeriod)
	} else {
		s.energyRecoveryPeriod = energyRecoveryPeriod
	}

	if minRewardsToClaim, err := conf.GetFloat64(context.Background(), "min_rewards_to_claim"); err != nil {
		log.Printf("[WARN] min_rewards_to_claim not found in config, using default value: %f", s.minRewardsToClaim)
	} else {
		s.minRewardsToClaim = minRewardsToClaim
	}

	if minVersion, err := conf.GetString(context.Background(), "min_version"); err != nil {
		log.Printf("[WARN] min_version not found in config, using default value: %s", s.minVersion)
	} else {
		s.minVersion = minVersion
	}

	if craftStepAmount, err := conf.GetFloat64(context.Background(), "craft_step_amount"); err != nil {
		log.Printf("[WARN] craft_step_amount not found in config, using default value: %f", s.craftStepAmount)
	} else {
		s.craftStepAmount = craftStepAmount
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

// GetCraftStepAmount ...
func (s *Service) GetCraftStepAmount(ctx context.Context) float64 {
	return s.craftStepAmount
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
	nfts, err := s.gameRepo.GetUserNFTs(ctx, uid)
	if err != nil {
		return nil, err
	}

	result := make([]NFTInfo, 0, len(nfts))
	for _, nft := range nfts {
		result = append(result, NFTInfo{
			ID:       nft.ID,
			MaxLevel: nft.MaxLevel,
			NftType:  NFTType(nft.NftType),
		})
	}

	return result, nil
}

// GetSelectedNFTID ...
func (s *Service) GetSelectedNFTID(ctx context.Context, uid uuid.UUID) (*string, error) {
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
func (s *Service) BuyNFTPack(ctx context.Context, uid, packID uuid.UUID) (*NFTInfo, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	pack, err := repo.GetNFTPack(ctx, packID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nft pack: %w", err)
	}

	nft, err := generateNFT(pack)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nft: %w", err)
	}

	res, err := repo.AddNFT(ctx, repository.AddNFTParams{
		UserID:   uid,
		ID:       nft.ID,
		NftType:  nft.NftType.String(),
		MaxLevel: int32(nft.MaxLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store nft: %w", err)
	}

	if err := repo.StoreSelectedNFT(ctx, repository.StoreSelectedNFTParams{
		UserID:        uid,
		SelectedNftID: sql.NullString{String: res.ID, Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to store selected nft: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return castDbNftInfoToNFTInfo(&res), nil
}

// CraftNFT ...
func (s *Service) CraftNFT(ctx context.Context, uid uuid.UUID, nftsToCraft []string) (*NFTInfo, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	nfts, err := repo.GetUserNFTByIDs(ctx, repository.GetUserNFTByIDsParams{
		UserID: uid,
		IDs:    nftsToCraft,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get nfts to craft: %w", err)
	}

	if len(nfts) != len(nftsToCraft) {
		return nil, ErrNotAllNftsToCraftWereFound
	}

	nft, err := craftNFT(nfts)
	if err != nil {
		return nil, fmt.Errorf("failed to craft nft: %w", err)
	}

	if _, err := repo.AddNFT(ctx, repository.AddNFTParams{
		UserID:   uid,
		ID:       nft.ID,
		NftType:  nft.NftType.String(),
		MaxLevel: int32(nft.MaxLevel),
	}); err != nil {
		return nil, fmt.Errorf("failed to store nft: %w", err)
	}

	if err := repo.CraftNFTs(ctx, repository.CraftNFTsParams{
		UserID:       uid,
		NftIds:       nftsToCraft,
		CraftedNftID: sql.NullString{String: nft.ID, Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to burn nfts: %w", err)
	}

	if err := repo.StoreSelectedNFT(ctx, repository.StoreSelectedNFTParams{
		UserID:        uid,
		SelectedNftID: sql.NullString{String: nft.ID, Valid: true},
	}); err != nil {
		return nil, fmt.Errorf("failed to store selected nft: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nft, nil
}

// SelectNFT ...
func (s *Service) SelectNFT(ctx context.Context, uid uuid.UUID, nftMintAddr string) error {
	nft, err := s.gameRepo.GetUserNFT(ctx, repository.GetUserNFTParams{
		UserID: uid,
		ID:     nftMintAddr,
	})
	if err != nil {
		return fmt.Errorf("failed to get user nft: %w", err)
	}

	if err := s.gameRepo.StoreSelectedNFT(ctx, repository.StoreSelectedNFTParams{
		UserID:        uid,
		SelectedNftID: sql.NullString{String: nft.ID, Valid: true},
	}); err != nil {
		return fmt.Errorf("failed to store selected nft: %w", err)
	}

	return nil
}

// StartGame ...
func (s *Service) StartGame(ctx context.Context, uid uuid.UUID, complexity int32, isTraining bool) error {
	player, err := s.GetPlayerInfo(ctx, uid)
	if err != nil {
		return err
	}

	nft, err := s.gameRepo.GetUserNFT(ctx, repository.GetUserNFTParams{
		UserID: uid,
		ID:     player.SelectedNftID,
	})
	if err != nil {
		return fmt.Errorf("failed to get user nft: %w", err)
	}

	if nft.MaxLevel < complexity {
		return fmt.Errorf("nft max level is %d, but complexity is %d", nft.MaxLevel, complexity)
	}

	if err := s.gameRepo.StartGame(ctx, repository.StartGameParams{
		UserID:     uid,
		NFTID:      player.SelectedNftID,
		Complexity: complexity,
		IsTraining: isTraining,
	}); err != nil {
		return fmt.Errorf("failed to start game: %w", err)
	}

	return nil
}

// FinishGame ...
// TODO: rewards calculation
func (s *Service) FinishGame(ctx context.Context, uid uuid.UUID, blocksDone int32) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	currentGame, err := repo.GetCurrentGame(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get current game: %w", err)
	}

	var rewardsAmount float64 = float64(blocksDone) * 0.5

	if err := repo.FinishGame(ctx, repository.FinishGameParams{
		ID:         currentGame.ID,
		BlocksDone: blocksDone,
	}); err != nil {
		return fmt.Errorf("failed to finish game: %w", err)
	}

	if err := repo.RewardsDeposit(ctx, repository.RewardsDepositParams{
		UserID:     uid,
		RelationID: uuid.NullUUID{UUID: currentGame.ID, Valid: true},
		Amount:     rewardsAmount,
	}); err != nil {
		return fmt.Errorf("failed to withdraw rewards: %w", err)
	}

	return tx.Commit()
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
	deposit, _ := s.gameRepo.GetUserRewardsDeposited(ctx, uid)
	withdrawn, _ := s.gameRepo.GetUserRewardsWithdrawn(ctx, uid)
	return math.Dim(deposit, withdrawn), nil
}

// claimRewardsFunc is a function that can be used to claim rewards to the user's SAO wallet
type claimRewardsFunc func(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error)

// ClaimRewards ...
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64, claimFn claimRewardsFunc) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	userRewardsAmount, _ := s.GetUserRewards(ctx, uid)
	if userRewardsAmount < s.GetMinAmountToClaim() {
		return fmt.Errorf("not enough rewards to claim, need %f, have %f", s.GetMinAmountToClaim(), userRewardsAmount)
	}
	if amount > userRewardsAmount {
		return fmt.Errorf("not enough rewards to claim, need %f, have %f", amount, userRewardsAmount)
	}

	if err := repo.RewardsWithdraw(ctx, repository.RewardsWithdrawParams{
		UserID: uid,
		Amount: amount,
	}); err != nil {
		return fmt.Errorf("failed to withdraw rewards: %w", err)
	}

	// Send tokens to user wallet
	if claimFn != nil {
		if _, err := claimFn(ctx, uid, userRewardsAmount); err != nil {
			return fmt.Errorf("failed to claim rewards: %w", err)
		}
	} else {
		log.Printf("no claim function provided, skipping claim")
	}

	return tx.Commit()
}

// Cast database result to NFTInfo struct
func castDbNftInfoToNFTInfo(dbNftInfo *repository.UnityGameNft) *NFTInfo {
	return &NFTInfo{
		ID:       dbNftInfo.ID,
		NftType:  NFTType(dbNftInfo.NftType),
		MaxLevel: dbNftInfo.MaxLevel,
	}
}
