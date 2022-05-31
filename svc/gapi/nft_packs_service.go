package gapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
	"github.com/google/uuid"
)

type (
	NFTPackService struct {
		repo nftPacksRepository
	}

	nftPacksRepository interface {
		AddNFTPack(ctx context.Context, arg repository.AddNFTPackParams) (repository.UnityGameNftPack, error)
		DeleteNFTPack(ctx context.Context, id uuid.UUID) error
		GetNFTPack(ctx context.Context, id uuid.UUID) (repository.UnityGameNftPack, error)
		GetNFTPacksList(ctx context.Context) ([]repository.UnityGameNftPack, error)
		SoftDeleteNFTPack(ctx context.Context, id uuid.UUID) error
		UpdateNFTPack(ctx context.Context, arg repository.UpdateNFTPackParams) (repository.UnityGameNftPack, error)
	}
)

// NewNFTPackService creates a new nft pack service
func NewNFTPackService(repo nftPacksRepository) *NFTPackService {
	return &NFTPackService{repo: repo}
}

// AddNFTPack adds a new nft pack
func (s *NFTPackService) AddNFTPack(ctx context.Context, name string, price float64, dropChances DropChances) (*NFTPackInfo, error) {
	nftPack, err := s.repo.AddNFTPack(ctx, repository.AddNFTPackParams{
		Name:        name,
		DropChances: dropChances.Bytes(),
		Price:       price,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new nft pack: %w", err)
	}

	return &NFTPackInfo{
		ID:          nftPack.ID.String(),
		Name:        nftPack.Name,
		DropChances: dropChances,
		Price:       nftPack.Price,
	}, nil
}

// DeleteNFTPack deletes a nft pack
func (s *NFTPackService) DeleteNFTPack(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteNFTPack(ctx, id); err != nil {
		return fmt.Errorf("failed to delete nft pack: %w", err)
	}

	return nil
}

// GetNFTPack gets a nft pack
func (s *NFTPackService) GetNFTPack(ctx context.Context, id uuid.UUID) (*NFTPackInfo, error) {
	nftPack, err := s.repo.GetNFTPack(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get nft pack: %w", err)
	}

	var dropChances DropChances
	if err := json.Unmarshal(nftPack.DropChances, &dropChances); err != nil {
		return nil, fmt.Errorf("failed to unmarshal drop chances: %w", err)
	}

	return &NFTPackInfo{
		ID:          nftPack.ID.String(),
		Name:        nftPack.Name,
		DropChances: dropChances,
		Price:       nftPack.Price,
	}, nil
}

// GetNFTPacksList gets a list of nft packs
func (s *NFTPackService) GetNFTPacksList(ctx context.Context) ([]NFTPackInfo, error) {
	nftPacks, err := s.repo.GetNFTPacksList(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nft packs list: %w", err)
	}

	var nftPacksList []NFTPackInfo
	for _, nftPack := range nftPacks {
		var dropChances DropChances
		if err := json.Unmarshal(nftPack.DropChances, &dropChances); err != nil {
			return nil, fmt.Errorf("failed to unmarshal drop chances: %w", err)
		}

		nftPacksList = append(nftPacksList, NFTPackInfo{
			ID:          nftPack.ID.String(),
			Name:        nftPack.Name,
			DropChances: dropChances,
			Price:       nftPack.Price,
		})
	}

	return nftPacksList, nil
}

// SoftDeleteNFTPack soft deletes a nft pack
func (s *NFTPackService) SoftDeleteNFTPack(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.SoftDeleteNFTPack(ctx, id); err != nil {
		return fmt.Errorf("failed to soft delete nft pack: %w", err)
	}

	return nil
}

// UpdateNFTPack updates a nft pack
func (s *NFTPackService) UpdateNFTPack(ctx context.Context, id uuid.UUID, name string, price float64, dropChances DropChances) (*NFTPackInfo, error) {
	nftPack, err := s.repo.UpdateNFTPack(ctx, repository.UpdateNFTPackParams{
		ID:          id,
		Name:        name,
		DropChances: dropChances.Bytes(),
		Price:       price,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update nft pack: %w", err)
	}

	return &NFTPackInfo{
		ID:          nftPack.ID.String(),
		Name:        nftPack.Name,
		DropChances: dropChances,
		Price:       nftPack.Price,
	}, nil
}
