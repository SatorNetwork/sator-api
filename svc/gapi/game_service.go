package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
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
		payment              paymentService
		minVersion           string
		energyFull           int32
		energyRecoveryPeriod time.Duration
		minRewardsToClaim    float64
		craftStepAmount      float64
		electricityMaxGames  int32
	}

	ServiceOption func(*Service)

	gameRepository interface {
		WithTx(tx *sql.Tx) *repository.Queries

		StartGame(ctx context.Context, arg repository.StartGameParams) error
		GetCurrentGame(ctx context.Context, userID uuid.UUID) (repository.UnityGameResult, error)
		FinishGame(ctx context.Context, arg repository.FinishGameParams) error
		SpendEnergyOfPlayer(ctx context.Context, userID uuid.UUID) error

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
		SelectCharacterToPlayer(ctx context.Context, arg repository.SelectCharacterToPlayerParams) error

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
		GetInt32(ctx context.Context, key string) (int32, error)
		GetJSON(ctx context.Context, key string, result interface{}) error
		GetDurration(ctx context.Context, key string) (time.Duration, error)
	}

	paymentService interface {
		GetBalance(ctx context.Context, uid uuid.UUID) (float64, error)
		ClaimRewards(ctx context.Context, uid uuid.UUID, amount, feePercent float64, feeDistribution map[string]float64) (string, error)
		Pay(ctx context.Context, uid uuid.UUID, amount float64, info string) (string, error)
	}

	PlayerInfo struct {
		UserID                uuid.UUID
		EnergyPoints          int
		EnergyPointsFull      int
		SelectedNftID         string
		ElectricityCost       float64
		ElectricitySpent      int32
		EnergyRecoveryPeriod  time.Duration
		EnergyRecoveryCurrent time.Duration
		ConversionFee         float64
		SelectedCharaterID    string
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(repo gameRepository, conf configer, p paymentService, opt ...ServiceOption) *Service {
	s := &Service{
		gameRepo:             repo,
		conf:                 conf,
		payment:              p,
		energyFull:           3,
		energyRecoveryPeriod: time.Hour * 4,
		minRewardsToClaim:    50,
		craftStepAmount:      500,
		electricityMaxGames:  18,
	}

	// Apply options
	for _, o := range opt {
		o(s)
	}

	return s
}

// GetPlayerInfo ...
func (s *Service) GetPlayerInfo(ctx context.Context, uid uuid.UUID) (*PlayerInfo, error) {
	energyFull, err := s.conf.GetInt32(ctx, "energy_full")
	if err != nil || energyFull == 0 {
		energyFull = s.energyFull
	}

	energyRecoveryPeriod, err := s.conf.GetDurration(context.Background(), "energy_recovery_period")
	if err != nil || energyRecoveryPeriod == 0 {
		energyRecoveryPeriod = s.energyRecoveryPeriod
	}

	player, err := s.gameRepo.GetPlayer(ctx, uid)
	if err != nil {
		if !db.IsNotFoundError(err) {
			return nil, err
		}

		player, err = s.gameRepo.AddNewPlayer(ctx, repository.AddNewPlayerParams{
			UserID:           uid,
			EnergyPoints:     energyFull,
			EnergyRefilledAt: time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	recoveryEnergy := recoveryEnergyPoints(player, energyFull, energyRecoveryPeriod)
	energy := recoveryEnergy + player.EnergyPoints

	lastEnergyRecoveredAt := player.EnergyRefilledAt
	if recoveryEnergy > 0 {
		lastEnergyRecoveredAt = player.EnergyRefilledAt.Add(energyRecoveryPeriod * time.Duration(recoveryEnergy))
		if err := s.gameRepo.RefillEnergyOfPlayer(ctx, repository.RefillEnergyOfPlayerParams{
			UserID:           uid,
			EnergyPoints:     energy,
			EnergyRefilledAt: lastEnergyRecoveredAt,
		}); err != nil {
			return nil, err
		}
	}

	timeToRecovery := time.Until(lastEnergyRecoveredAt.Add(energyRecoveryPeriod))

	convFee, err := s.conf.GetFloat64(ctx, "convert_commission")
	if err != nil {
		convFee = 0
	}

	return &PlayerInfo{
		UserID:                player.UserID,
		EnergyPoints:          int(energy),
		EnergyPointsFull:      int(energyFull),
		SelectedNftID:         player.SelectedNftID.String,
		ElectricityCost:       player.ElectricityCosts,
		ElectricitySpent:      player.ElectricitySpent,
		EnergyRecoveryPeriod:  energyRecoveryPeriod,
		EnergyRecoveryCurrent: timeToRecovery,
		ConversionFee:         convFee,
		SelectedCharaterID:    player.SelectedCharacterID.String,
	}, nil
}

// GetUserBalance ...
func (s *Service) GetUserBalance(ctx context.Context, uid uuid.UUID) (float64, error) {
	return s.payment.GetBalance(ctx, uid)
}

// GetMinVersion ...
func (s *Service) GetMinVersion(ctx context.Context) string {
	minVersion, err := s.conf.GetString(context.Background(), "min_version")
	if err != nil || minVersion == "" {
		minVersion = s.minVersion
	}

	return minVersion
}

// GetCraftStepAmount ...
func (s *Service) GetCraftStepAmount(ctx context.Context) float64 {
	craftStepAmount, err := s.conf.GetFloat64(context.Background(), "craft_step_amount")
	if err != nil || craftStepAmount == 0 {
		craftStepAmount = s.craftStepAmount
	}

	return craftStepAmount
}

func (s *Service) GetElectricityMaxGames(ctx context.Context) int32 {
	electricityMaxGames, err := s.conf.GetInt(context.Background(), "electricity_max_games")
	if err != nil || electricityMaxGames == 0 {
		return s.electricityMaxGames
	}

	return int32(electricityMaxGames)
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

	if pack.Price > 0 {
		userBalance, _ := s.GetUserBalance(ctx, uid)
		if userBalance < pack.Price {
			return nil, ErrInsufficientBalance
		}

		if tr, err := s.payment.Pay(ctx, uid, pack.Price, "purchase of nft pack"); err != nil {
			log.Printf("failed to buy nft pack: %v", err)
			return nil, fmt.Errorf("failed to buy nft pack: %w", err)
		} else {
			log.Printf("successful purchase of nft pack: %s", tr)
		}
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

	craftStepAmount := s.GetCraftStepAmount(ctx)
	userBalance, _ := s.GetUserBalance(ctx, uid)

	nextNftType := getNextNFTType(NFTType(nfts[0].NftType))
	if nextNftType.ToInt() < 1 {
		return nil, ErrCouldNotCraftNFT
	}

	craftCost := float64(nextNftType.ToInt()) * craftStepAmount
	if userBalance < craftCost {
		return nil, fmt.Errorf("not enough balance to craft nft: you have %f but you need %f", userBalance, craftCost)
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

	if craftCost > 0 {
		if tr, err := s.payment.Pay(ctx, uid, craftCost, "crafting in-game nft"); err != nil {
			log.Printf("failed to pay for crafting nft: %v", err)
			return nil, ErrCouldNotCraftNFT
		} else {
			log.Printf("successful payment for crafting: %s", tr)
		}
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

func (s *Service) SelectCharacter(ctx context.Context, userID uuid.UUID, characterID string) error {
	if err := s.gameRepo.SelectCharacterToPlayer(ctx, repository.SelectCharacterToPlayerParams{
		SelectedCharacterID: sql.NullString{
			String: characterID,
			Valid:  true,
		},
		UserID: userID,
	}); err != nil {
		return fmt.Errorf("failed to select character: %w", err)
	}

	return nil
}

// StartGame ...
func (s *Service) StartGame(ctx context.Context, uid uuid.UUID, complexity int32, isTraining bool) (*GameConfig, error) {
	log.Printf("start game: %s, %d, %t", uid, complexity, isTraining)

	leftElectr, _, _ := s.GetElectricityLeft(ctx, uid)
	if leftElectr < 1 {
		log.Printf("not enough electricity to start game")
		return nil, ErrNotEnoughElectricity
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	player, err := s.GetPlayerInfo(ctx, uid)
	if err != nil {
		return nil, err
	}

	nft, err := repo.GetUserNFT(ctx, repository.GetUserNFTParams{
		UserID: uid,
		ID:     player.SelectedNftID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user nft: %w", err)
	}

	if nft.MaxLevel < complexity {
		return nil, fmt.Errorf("nft max level is %d, but complexity is %d", nft.MaxLevel, complexity)
	}

	gameConfig := &GameConfig{}
	confName := fmt.Sprintf("complexity_%s", getGameLevelName(complexity))
	if err := s.conf.GetJSON(ctx, confName, gameConfig); err != nil {
		return nil, fmt.Errorf("failed to get game config %s: %w", confName, err)
	}

	if err := repo.StartGame(ctx, repository.StartGameParams{
		UserID:     uid,
		NFTID:      player.SelectedNftID,
		Complexity: complexity,
		IsTraining: isTraining,
	}); err != nil {
		return nil, fmt.Errorf("failed to start game: %w", err)
	}

	if err := repo.SpendEnergyOfPlayer(ctx, uid); err != nil {
		return nil, fmt.Errorf("failed to take the energy of player: %w", err)
	}

	if player.EnergyPoints == player.EnergyPointsFull {
		if err := repo.ResetEnergyRefilledAtOfPlayer(ctx, uid); err != nil {
			return nil, fmt.Errorf("failed to update energy refilled at of player: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return gameConfig, nil
}

// FinishGame ...
// TODO: rewards calculation
func (s *Service) FinishGame(ctx context.Context, uid uuid.UUID, result, blocksDone int32) (int, error) {
	log.Printf("finish game: %s, %d, %d", uid, result, blocksDone)

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	currentGame, err := repo.GetCurrentGame(ctx, uid)
	if err != nil {
		return 0, fmt.Errorf("failed to get current game: %w", err)
	}

	var (
		rewardsAmount, electricityCost float64
		electricitySpent               int32
		viewers                        int
	)

	if !currentGame.IsTraining {
		nft, err := repo.GetUserNFT(ctx, repository.GetUserNFTParams{
			UserID: uid,
			ID:     currentGame.NFTID,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to get current nft: %w", err)
		}

		rewardsAmount, viewers, err = calculateUserRewardsForGame(s.conf, nft.NftType, currentGame.Complexity, result)
		if err != nil {
			return 0, fmt.Errorf("failed to calculate user rewards: %w", err)
		}

		electricityCost, err = calculateElectricityCost(s.conf, nft.NftType, result, rewardsAmount)
		if err != nil {
			return 0, fmt.Errorf("failed to calculate electricity cost: %w", err)
		}

		if electricityCost > 0 {
			electricitySpent = 1
		}
	}

	log.Printf("rewards amount: %f, electricity cost: %f", rewardsAmount, electricityCost)

	if err := repo.FinishGame(ctx, repository.FinishGameParams{
		ID:               currentGame.ID,
		BlocksDone:       blocksDone,
		Result:           sql.NullInt32{Int32: result, Valid: true},
		ElectricityCosts: electricityCost,
	}); err != nil {
		return 0, fmt.Errorf("failed to finish game: %w", err)
	}

	if err := repo.RewardsDeposit(ctx, repository.RewardsDepositParams{
		UserID:     uid,
		RelationID: uuid.NullUUID{UUID: currentGame.ID, Valid: true},
		Amount:     rewardsAmount,
	}); err != nil {
		return 0, fmt.Errorf("failed to withdraw rewards: %w", err)
	}

	if err := repo.AddElectricityToPlayer(ctx, repository.AddElectricityToPlayerParams{
		UserID:           uid,
		ElectricityCosts: electricityCost,
		ElectricitySpent: electricitySpent,
	}); err != nil {
		return 0, fmt.Errorf("failed to take the energy of player: %w", err)
	}

	return viewers, tx.Commit()
}

// GetMinAmountToClaim ...
func (s *Service) GetMinAmountToClaim() float64 {
	minRewardsToClaim, err := s.conf.GetFloat64(context.Background(), "min_rewards_to_claim")
	if err != nil || minRewardsToClaim == 0 {
		minRewardsToClaim = s.minRewardsToClaim
	}

	return minRewardsToClaim
}

// GetUserRewards ...
func (s *Service) GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error) {
	deposit, err := s.gameRepo.GetUserRewardsDeposited(ctx, uid)
	if err != nil {
		log.Printf("failed to get user rewards deposited: %v", err)
		return 0, nil
	}

	withdrawn, err := s.gameRepo.GetUserRewardsWithdrawn(ctx, uid)
	if err != nil {
		log.Printf("failed to get user rewards withdrawn: %v", err)
		return deposit, nil
	}

	result, _ := big.NewFloat(0).Sub(big.NewFloat(deposit), big.NewFloat(withdrawn)).Float64()
	log.Printf("user %s rewards: %f", uid, result)
	return result, nil
}

// ClaimRewards ...
func (s *Service) ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	minAmountToClaim := s.GetMinAmountToClaim()
	if amount < minAmountToClaim {
		return fmt.Errorf("amount to claim is less than %f", minAmountToClaim)
	}

	userRewardsAmount, _ := s.GetUserRewards(ctx, uid)
	if userRewardsAmount < minAmountToClaim {
		return fmt.Errorf("not enough rewards to claim, need %f, have %f", s.GetMinAmountToClaim(), userRewardsAmount)
	}
	if amount > userRewardsAmount {
		return fmt.Errorf("not enough rewards to claim, you want %f, but have %f", amount, userRewardsAmount)
	}

	if err := repo.RewardsWithdraw(ctx, repository.RewardsWithdrawParams{
		UserID: uid,
		Amount: amount,
	}); err != nil {
		return fmt.Errorf("failed to withdraw rewards: %w", err)
	}

	fee, _ := s.conf.GetFloat64(ctx, "convert_commission")
	feeDistr := make(map[string]float64)
	if fee > 0 {
		if err := s.conf.GetJSON(ctx, "conversion_fee_accumulator", &feeDistr); err != nil {
			return fmt.Errorf("failed to get convert fee distribution: %w", err)
		}
	}

	if tr, err := s.payment.ClaimRewards(ctx, uid, amount, fee, feeDistr); err != nil {
		log.Printf("failed to claim rewards: %v", err)
		return ErrCouldNotClaimRewards
	} else {
		log.Printf("successful claim rewards: %s", tr)
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

func (s *Service) GetElectricityLeft(ctx context.Context, uid uuid.UUID) (left, max int32, err error) {
	electricityMax := s.GetElectricityMaxGames(ctx)
	player, err := s.gameRepo.GetPlayer(ctx, uid)
	if err != nil {
		return 0, electricityMax, fmt.Errorf("failed to get player: %w", err)
	}

	return electricityMax - player.ElectricitySpent, electricityMax, nil
}

func (s *Service) PayForElectricity(ctx context.Context, uid uuid.UUID) error {
	log.Printf("pay for electricity")

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	repo := s.gameRepo.WithTx(tx)

	player, err := repo.GetPlayer(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get player: %w", err)
	}

	log.Printf("PayForElectricity: player: %+v", player)

	if player.ElectricityCosts <= 0 {
		return nil
	}

	balance, err := s.GetUserBalance(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to get user balance: %w", err)
	}

	log.Printf("PayForElectricity: balance: %+v", balance)

	if balance < player.ElectricityCosts {
		log.Printf("PayForElectricity: not enough balance; balance: %v, need to pay: %v", balance, player.ElectricityCosts)
		return fmt.Errorf("not enough funds to pay for electricity")
	}

	if err := repo.ResetElectricityForPlayer(ctx, uid); err != nil {
		log.Printf("failed to reset electricity for player: %v", err)
		return fmt.Errorf("failed to reset electricity for player: %w", err)
	}

	if tr, err := s.payment.Pay(ctx, uid, player.ElectricityCosts, "electricity"); err != nil {
		log.Printf("failed to pay for electricity: %v", err)
		return ErrCouldNotPayForElectricity
	} else {
		log.Printf("successful payment for electricity: %s", tr)
	}

	return tx.Commit()
}
