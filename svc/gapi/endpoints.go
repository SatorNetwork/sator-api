package gapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetStatus    endpoint.Endpoint
		GetNFTPacks  endpoint.Endpoint
		BuyNFTPack   endpoint.Endpoint
		CraftNFT     endpoint.Endpoint
		SelectNFT    endpoint.Endpoint
		StartGame    endpoint.Endpoint
		FinishGame   endpoint.Endpoint
		ClaimRewards endpoint.Endpoint
	}

	gameService interface {
		GetPlayerInfo(ctx context.Context, uid uuid.UUID) (*PlayerInfo, error)
		GetMinVersion(ctx context.Context) string

		GetUserNFTs(ctx context.Context, uid uuid.UUID) ([]NFTInfo, error)
		GetNFTPacks(ctx context.Context, uid uuid.UUID) ([]NFTPackInfo, error)
		BuyNFTPack(ctx context.Context, uid, packID uuid.UUID) (*NFTInfo, error)
		CraftNFT(ctx context.Context, uid uuid.UUID, nftsToCraft []string) (*NFTInfo, error)
		SelectNFT(ctx context.Context, uid uuid.UUID, nftMintAddr string) error

		StartGame(ctx context.Context, uid uuid.UUID, complexity int32, isTraining bool) error
		FinishGame(ctx context.Context, uid uuid.UUID, blocksDone int32) error

		GetDefaultGameConfig(ctx context.Context, uid uuid.UUID) (*GameConfig, error)

		GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error)
		ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64, claimFn claimRewardsFunc) error
		GetMinAmountToClaim() float64
	}

	walletService interface {
		GetUserBalance(ctx context.Context, uid uuid.UUID) (float64, error)
		ClaimInGameRewards(ctx context.Context, userID uuid.UUID, amount float64) (tx string, err error)
	}
)

func MakeEndpoints(
	gs gameService,
	ws walletService,
	m ...endpoint.Middleware,
) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetStatus:    MakeGetStatusEndpoint(gs, ws),
		GetNFTPacks:  MakeGetNFTPacksEndpoint(gs),
		BuyNFTPack:   MakeBuyNFTPackEndpoint(gs, validateFunc),
		CraftNFT:     MakeCraftNFTEndpoint(gs, validateFunc),
		SelectNFT:    MakeSelectNFTEndpoint(gs, validateFunc),
		StartGame:    MakeStartGameEndpoint(gs, validateFunc),
		FinishGame:   MakeFinishGameEndpoint(gs, validateFunc),
		ClaimRewards: MakeClaimRewardsEndpoint(gs, ws, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetStatus = mdw(e.GetStatus)
			e.GetNFTPacks = mdw(e.GetNFTPacks)
			e.BuyNFTPack = mdw(e.BuyNFTPack)
			e.CraftNFT = mdw(e.CraftNFT)
			e.SelectNFT = mdw(e.SelectNFT)
			e.StartGame = mdw(e.StartGame)
			e.FinishGame = mdw(e.FinishGame)
			e.ClaimRewards = mdw(e.ClaimRewards)
		}
	}

	return e
}

type GetStatusResponse struct {
	EnergyLeft                   int       `json:"energy_left"`
	UserCurrency                 float64   `json:"user_currency"`
	UserInGameCurrency           float64   `json:"user_in_game_currency"`
	MinAmountOfCurrencyToConvert float64   `json:"min_amount_of_currency_to_convert"`
	MinVersion                   string    `json:"min_version"`
	SelectedNFTID                *string   `json:"selected_nft_id"`
	UserOwnedNFTList             []NFTInfo `json:"user_owned_nft_list"`
}

// MakeGetStatusEndpoint ...
func MakeGetStatusEndpoint(s gameService, ws walletService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		player, err := s.GetPlayerInfo(ctx, uid)
		if err != nil {
			return nil, err
		}

		totalRewards, err := s.GetUserRewards(ctx, uid)
		if err != nil {
			log.Printf("could not get user rewards: %v", err)
		}

		userCurrency, err := ws.GetUserBalance(ctx, uid)
		if err != nil {
			log.Printf("could not get user balance: %v", err)
		}

		userNFTs, err := s.GetUserNFTs(ctx, uid)
		if err != nil {
			log.Printf("could not get user nfts: %v", err)
		}

		var selectedNFT *string
		if player.SelectedNftID != "" {
			selectedNFT = &player.SelectedNftID
		}

		return GetStatusResponse{
			EnergyLeft:                   player.EnergyPoints,
			UserCurrency:                 userCurrency,
			UserInGameCurrency:           totalRewards,
			MinAmountOfCurrencyToConvert: s.GetMinAmountToClaim(),
			MinVersion:                   s.GetMinVersion(ctx),
			SelectedNFTID:                selectedNFT,
			UserOwnedNFTList:             userNFTs,
		}, nil
	}
}

// MakeGetNFTPacksEndpoint ...
func MakeGetNFTPacksEndpoint(s gameService) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		packs, err := s.GetNFTPacks(ctx, uid)
		if err != nil {
			return nil, err
		}

		return packs, nil
	}
}

type (
	BuyNFTPackRequest struct {
		PackID string `json:"pack_id" validate:"required"`
	}

	BuyNFTPackResponse struct {
		NewNFT           *NFTInfo  `json:"new_nft"`
		UserOwnedNftList []NFTInfo `json:"user_owned_nft_list"`
	}
)

// MakeBuyNFTPackEndpoint ...
func MakeBuyNFTPackEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(BuyNFTPackRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		pid, err := uuid.Parse(req.PackID)
		if err != nil {
			return nil, fmt.Errorf("could not parse pack id: %w", err)
		}

		newNFT, err := s.BuyNFTPack(ctx, uid, pid)
		if err != nil {
			return nil, err
		}

		userNFTs, err := s.GetUserNFTs(ctx, uid)
		if err != nil {
			return nil, err
		}

		return BuyNFTPackResponse{
			NewNFT:           newNFT,
			UserOwnedNftList: userNFTs,
		}, nil
	}
}

type (
	CraftNFTRequest struct {
		NFTsToCraft []string `json:"nfts_to_craft" validate:"required,min=1"`
	}

	CraftNFTResponse struct {
		NewNFT           *NFTInfo  `json:"new_nft"`
		UserOwnedNFTList []NFTInfo `json:"user_owned_nft_list"`
	}
)

// MakeCraftNFTEndpoint ...
func MakeCraftNFTEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(CraftNFTRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		newNFT, err := s.CraftNFT(ctx, uid, req.NFTsToCraft)
		if err != nil {
			return nil, err
		}

		userNFTs, err := s.GetUserNFTs(ctx, uid)
		if err != nil {
			return nil, err
		}

		return CraftNFTResponse{
			NewNFT:           newNFT,
			UserOwnedNFTList: userNFTs,
		}, nil
	}
}

type (
	SelectNFTRequest struct {
		NFTID string `json:"nft_id" validate:"required"`
	}

	SelectNFTResponse struct {
		UserOwnedNFTList []NFTInfo `json:"user_owned_nft_list"`
	}
)

// MakeSelectNFTEndpoint ...
func MakeSelectNFTEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(SelectNFTRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.SelectNFT(ctx, uid, req.NFTID); err != nil {
			return nil, err
		}

		return true, nil
	}
}

type (
	StartGameRequest struct {
		SelectedComplexity int32 `json:"selected_complexity" validate:"required"`
		IsTraining         bool  `json:"is_training"`
	}

	StartGameResponse struct {
		GameConfigJSON string `json:"game_config_json"`
	}
)

// MakeStartGameEndpoint ...
func MakeStartGameEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(StartGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.StartGame(ctx, uid, req.SelectedComplexity, req.IsTraining); err != nil {
			return nil, err
		}

		gameConfig, err := s.GetDefaultGameConfig(ctx, uid)
		if err != nil {
			return nil, err
		}

		config, err := json.Marshal(gameConfig)
		if err != nil {
			return nil, err
		}

		return StartGameResponse{
			GameConfigJSON: string(config),
		}, nil
	}
}

type (
	FinishGameRequest struct {
		BlocksDone int32 `json:"blocks_done" validate:"required"`
	}

	FinishGameResponse struct {
		UserInGameCurrency float64 `json:"user_in_game_currency"`
	}
)

// MakeFinishGameEndpoint ...
func MakeFinishGameEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(FinishGameRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.FinishGame(ctx, uid, req.BlocksDone); err != nil {
			return nil, err
		}

		rewardsAmount, err := s.GetUserRewards(ctx, uid)
		if err != nil {
			return nil, err
		}

		return FinishGameResponse{
			UserInGameCurrency: rewardsAmount,
		}, nil
	}
}

// ClaimRewardsRequest ...
type ClaimRewardsRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// MakeClaimRewardsEndpoint ...
func MakeClaimRewardsEndpoint(s gameService, ws walletService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(ClaimRewardsRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.ClaimRewards(ctx, uid, req.Amount, ws.ClaimInGameRewards); err != nil {
			return nil, err
		}

		return true, nil
	}
}
