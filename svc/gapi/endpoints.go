package gapi

import (
	"context"
	"fmt"
	"log"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		GetStatus         endpoint.Endpoint
		GetNFTPacks       endpoint.Endpoint
		BuyNFTPack        endpoint.Endpoint
		CraftNFT          endpoint.Endpoint
		SelectNFT         endpoint.Endpoint
		SelectCharacter   endpoint.Endpoint
		StartGame         endpoint.Endpoint
		FinishGame        endpoint.Endpoint
		ClaimRewards      endpoint.Endpoint
		PayForElectricity endpoint.Endpoint

		GetSettingsValueTypes endpoint.Endpoint
		GetSettings           endpoint.Endpoint
		GetSettingsByKey      endpoint.Endpoint
		AddSetting            endpoint.Endpoint
		UpdateSetting         endpoint.Endpoint
		DeleteSetting         endpoint.Endpoint
	}

	gameService interface {
		GetPlayerInfo(ctx context.Context, uid uuid.UUID) (*PlayerInfo, error)
		GetMinVersion(ctx context.Context) string
		GetCraftStepAmount(ctx context.Context) float64

		GetElectricityLeft(ctx context.Context, uid uuid.UUID) (left, max int32, err error)
		PayForElectricity(ctx context.Context, uid uuid.UUID) error

		GetUserNFTs(ctx context.Context, uid uuid.UUID) ([]NFTInfo, error)
		GetNFTPacks(ctx context.Context, uid uuid.UUID) ([]NFTPackInfo, error)
		BuyNFTPack(ctx context.Context, uid, packID uuid.UUID) (*NFTInfo, error)
		CraftNFT(ctx context.Context, uid uuid.UUID, nftsToCraft []string) (*NFTInfo, error)
		SelectNFT(ctx context.Context, uid uuid.UUID, nftMintAddr string) error
		SelectCharacter(ctx context.Context, userID uuid.UUID, characterID string) error

		StartGame(ctx context.Context, uid uuid.UUID, complexity int32, isTraining bool) (*GameConfig, error)
		FinishGame(ctx context.Context, uid uuid.UUID, gameResult, blocksDone int32) (int, error)

		GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error)
		ClaimRewards(ctx context.Context, uid uuid.UUID, amount float64) error
		GetMinAmountToClaim() float64
		GetUserBalance(ctx context.Context, uid uuid.UUID) (float64, error)
	}

	gameSettingsService interface {
		Add(ctx context.Context, key, name, valueType string, value interface{}, description string) (Settings, error)
		Get(ctx context.Context, key string) (Settings, error)
		GetAll(ctx context.Context) []Settings
		Update(ctx context.Context, key string, value interface{}) (Settings, error)
		Delete(ctx context.Context, key string) error
		SettingsValueTypes() map[string]string
	}
)

func MakeEndpoints(
	gs gameService,
	settings gameSettingsService,
	m ...endpoint.Middleware,
) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetStatus:         MakeGetStatusEndpoint(gs, settings),
		GetNFTPacks:       MakeGetNFTPacksEndpoint(gs),
		BuyNFTPack:        MakeBuyNFTPackEndpoint(gs, validateFunc),
		CraftNFT:          MakeCraftNFTEndpoint(gs, validateFunc),
		SelectNFT:         MakeSelectNFTEndpoint(gs, validateFunc),
		StartGame:         MakeStartGameEndpoint(gs, validateFunc),
		FinishGame:        MakeFinishGameEndpoint(gs, validateFunc),
		ClaimRewards:      MakeClaimRewardsEndpoint(gs, validateFunc),
		PayForElectricity: MakePayForElectricityEndpoint(gs),
		SelectCharacter:   MakeSelectCharacterEndpoint(gs, validateFunc),

		GetSettings:           MakeGetSettingsEndpoint(settings),
		GetSettingsByKey:      MakeGetSettingsByKeyEndpoint(settings),
		AddSetting:            MakeAddSettingEndpoint(settings, validateFunc),
		UpdateSetting:         MakeUpdateSettingEndpoint(settings, validateFunc),
		DeleteSetting:         MakeDeleteSettingEndpoint(settings),
		GetSettingsValueTypes: MakeGetSettingsValueTypesEndpoint(settings),
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
			e.PayForElectricity = mdw(e.PayForElectricity)
			e.SelectCharacter = mdw(e.SelectCharacter)

			e.GetSettings = mdw(e.GetSettings)
			e.GetSettingsByKey = mdw(e.GetSettingsByKey)
			e.AddSetting = mdw(e.AddSetting)
			e.UpdateSetting = mdw(e.UpdateSetting)
			e.DeleteSetting = mdw(e.DeleteSetting)
			e.GetSettingsValueTypes = mdw(e.GetSettingsValueTypes)
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
	SelectedCharaterID           string    `json:"selected_character_id"`
	UserOwnedNFTList             []NFTInfo `json:"user_owned_nft_list"`
	CraftStepAmount              float64   `json:"craft_step_amount"`
	ElectricityLeft              int32     `json:"electricity_left"`
	ElectricityCost              float64   `json:"electricity_cost"`
	ElectricityMaxGames          int32     `json:"electricity_max_games"`
	EnergyRecoveryPeriod         int64     `json:"energy_recovery_period"`
	EnergyRecoveryCurrent        int64     `json:"energy_recovery_current"`
	ConversionFee                float64   `json:"convert_commission"`
}

// MakeGetStatusEndpoint ...
func MakeGetStatusEndpoint(s gameService, settings gameSettingsService) endpoint.Endpoint {
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

		userCurrency, err := s.GetUserBalance(ctx, uid)
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

		electrLeft, electrMax, _ := s.GetElectricityLeft(ctx, uid)

		resp := GetStatusResponse{
			EnergyLeft:                   player.EnergyPoints,
			UserCurrency:                 userCurrency,
			UserInGameCurrency:           totalRewards,
			MinAmountOfCurrencyToConvert: s.GetMinAmountToClaim(),
			MinVersion:                   s.GetMinVersion(ctx),
			SelectedNFTID:                selectedNFT,
			SelectedCharaterID:           player.SelectedCharaterID,
			UserOwnedNFTList:             userNFTs,
			CraftStepAmount:              s.GetCraftStepAmount(ctx),
			ElectricityLeft:              electrLeft,
			ElectricityCost:              player.ElectricityCost,
			ElectricityMaxGames:          electrMax,
			EnergyRecoveryPeriod:         int64(player.EnergyRecoveryPeriod.Seconds()),
			EnergyRecoveryCurrent:        int64(player.EnergyRecoveryCurrent.Seconds()),
			ConversionFee:                player.ConversionFee,
		}

		return resp, nil
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

		resp := BuyNFTPackResponse{
			NewNFT:           newNFT,
			UserOwnedNftList: userNFTs,
		}

		return resp, nil
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

		resp := CraftNFTResponse{
			NewNFT:           newNFT,
			UserOwnedNFTList: userNFTs,
		}

		return resp, nil
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
		GameConfig *GameConfig `json:"game_config"`
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

		gameConfig, err := s.StartGame(ctx, uid, req.SelectedComplexity, req.IsTraining)
		if err != nil {
			return nil, err
		}

		resp := StartGameResponse{
			GameConfig: gameConfig,
		}

		return resp, nil
	}
}

type (
	FinishGameRequest struct {
		BlocksDone int32 `json:"blocks_done" validate:"gte=0"`
		GameResult int32 `json:"game_result" validate:"oneof=0 1"`
	}

	FinishGameResponse struct {
		UserInGameCurrency float64 `json:"user_in_game_currency"`
		Viewers            int     `json:"viewers"`
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

		viewers, err := s.FinishGame(ctx, uid, req.GameResult, req.BlocksDone)
		if err != nil {
			return nil, err
		}

		rewardsAmount, err := s.GetUserRewards(ctx, uid)
		if err != nil {
			return nil, err
		}

		resp := FinishGameResponse{
			UserInGameCurrency: rewardsAmount,
			Viewers:            viewers,
		}

		return resp, nil
	}
}

// ClaimRewardsRequest ...
type ClaimRewardsRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

// MakeClaimRewardsEndpoint ...
func MakeClaimRewardsEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(ClaimRewardsRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.ClaimRewards(ctx, uid, req.Amount); err != nil {
			return nil, err
		}

		return true, nil
	}
}

type PayForElectricityResponse struct {
	ElectricityLeft int32 `json:"electricity_left"`
}

// MakePayForElectricityEndpoint ...
func MakePayForElectricityEndpoint(s gameService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		if err := s.PayForElectricity(ctx, uid); err != nil {
			return nil, err
		}

		left, _, err := s.GetElectricityLeft(ctx, uid)
		if err != nil {
			return nil, err
		}

		resp := PayForElectricityResponse{
			ElectricityLeft: left,
		}

		return resp, nil
	}
}

type (
	SelectCharacterRequest struct {
		CharacterID string `json:"character_id" validate:"required"`
	}
)

func MakeSelectCharacterEndpoint(s gameService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(SelectCharacterRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		if err := s.SelectCharacter(ctx, userID, req.CharacterID); err != nil {
			return false, err
		}

		return true, nil
	}
}

// MakeGetSettingsEndpoint ...
func MakeGetSettingsEndpoint(s gameSettingsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		return s.GetAll(ctx), nil
	}
}

// AddGameSettingsRequest ...
type AddGameSettingsRequest struct {
	Key         string      `json:"key" validate:"required"`
	Name        string      `json:"name"`
	ValueType   string      `json:"value_type" validate:"required,oneof=int float string json bool duration datetime"`
	Value       interface{} `json:"value" validate:"required"`
	Description string      `json:"description,omitempty"`
}

// MakeAddSettingEndpoint ...
func MakeAddSettingEndpoint(s gameSettingsService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		req := request.(AddGameSettingsRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.Add(ctx, req.Key, req.Name, req.ValueType, req.Value, req.Description)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// UpdateGameSettingRequest ...
type UpdateGameSettingRequest struct {
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

// MakeUpdateSettingEndpoint ...
func MakeUpdateSettingEndpoint(s gameSettingsService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		req := request.(UpdateGameSettingRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.Update(ctx, req.Key, req.Value)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// MakeDeleteSettingEndpoint ...
func MakeDeleteSettingEndpoint(s gameSettingsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		if err := s.Delete(ctx, request.(string)); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetSettingsValueTypesEndpoint ...
func MakeGetSettingsValueTypesEndpoint(s gameSettingsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		return s.SettingsValueTypes(), nil
	}
}

// MakeGetSettingsByKeyEndpoint ...
func MakeGetSettingsByKeyEndpoint(s gameSettingsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		return s.Get(ctx, request.(string))
	}
}
