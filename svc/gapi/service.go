package gapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		minVersion string
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService() *Service {
	return &Service{}
}

// GetMinVersion ...
func (s *Service) GetMinVersion(ctx context.Context) string {
	return s.minVersion
}

// GetEnergyLeft ...
func (s *Service) GetEnergyLeft(ctx context.Context, uid uuid.UUID) (int, error) {
	return 1, nil
}

// GetUserNFTs ...
func (s *Service) GetUserNFTs(ctx context.Context, uid uuid.UUID) ([]NFTInfo, error) {
	return nil, nil
}

// GetSelectedNFT ...
func (s *Service) GetSelectedNFT(ctx context.Context, uid uuid.UUID) (*string, error) {
	return utils.StringPointer(uuid.New().String()), nil
}

// GetNFTPacks ...
func (s *Service) GetNFTPacks(ctx context.Context, uid uuid.UUID) ([]NFTPackInfo, error) {
	return nil, nil
}

// BuyNFTPack ...
func (s *Service) BuyNFTPack(ctx context.Context, uid, packID uuid.UUID) (bool, error) {
	return false, nil
}

// CraftNFT ...
func (s *Service) CraftNFT(ctx context.Context, uid uuid.UUID, nftsToCraft []string) (bool, error) {
	return false, nil
}

// SelectNFT ...
func (s *Service) SelectNFT(ctx context.Context, uid uuid.UUID, nftMintAddr string) (bool, error) {
	return false, nil
}

// StartGame ...
func (s *Service) StartGame(ctx context.Context, uid uuid.UUID, complexity string, isTraining bool) (bool, error) {
	return true, nil
}

// FinishGame ...
func (s *Service) FinishGame(ctx context.Context, uid uuid.UUID, blocksDone int) (bool, error) {
	return true, nil
}

// GetDefaultGameConfig ...
func (s *Service) GetDefaultGameConfig(ctx context.Context, uid uuid.UUID) (*GameConfig, error) {
	cnf := &GameConfig{}
	if err := json.Unmarshal([]byte(defaultGameConfig), cnf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal game config: %w", err)
	}

	return cnf, nil
}
