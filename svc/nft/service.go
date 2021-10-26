package nft

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	"github.com/SatorNetwork/sator-api/internal/db"
	"github.com/SatorNetwork/sator-api/svc/nft/repository"
	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct {
		nftRepo    nftRepository
		buyNFTFunc buyNFTFunction
	}

	NFT struct {
		ID          uuid.UUID
		ImageLink   string
		Name        string
		Description string
		Tags        map[string]string
		// Supply - the number of copies that can be minted.
		Supply int
		// Royalties are optional and allow user to earn a percentage on secondary sales
		Royalties  float64 // TODO(evg): add validation?
		Blockchain string  // TODO(evg): replace with enum?
		SellType   string  // TODO(evg): replace with enum?

		BuyNowPrice   float64
		AuctionParams *NFTAuctionParams

		// NFT payload, e.g.: link to the original file, etc
		TokenURI string
	}

	NFTAuctionParams struct {
		StartingBid    float64
		StartTimestamp string // TODO(evg): replace with time.Time?
		EndTimestamp   string // TODO(evg): replace with time.Time?
	}

	Category struct {
		ID    uuid.UUID
		Title string
	}

	// Option func to set custom service options
	Option func(*Service)

	nftRepository interface {
		AddNFTItem(ctx context.Context, arg repository.AddNFTItemParams) (repository.NFTItem, error)
		GetNFTItemByID(ctx context.Context, id uuid.UUID) (repository.NFTItem, error)
		GetNFTItemsList(ctx context.Context, arg repository.GetNFTItemsListParams) ([]repository.NFTItem, error)
		GetNFTItemsListByRelationID(ctx context.Context, arg repository.GetNFTItemsListByRelationIDParams) ([]repository.NFTItem, error)
		GetNFTItemsListByOwnerID(ctx context.Context, arg repository.GetNFTItemsListByOwnerIDParams) ([]repository.NFTItem, error)
		GetNFTCategoriesList(ctx context.Context) ([]repository.NFTCategory, error)
		GetMainNFTCategory(ctx context.Context) (repository.NFTCategory, error)
		UpdateNFTItemOwner(ctx context.Context, arg repository.UpdateNFTItemOwnerParams) error
	}

	// Simple function
	buyNFTFunction func(ctx context.Context, uid uuid.UUID, amount float64, info string) error
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(nftRepo nftRepository, buyNFTFunc buyNFTFunction, opt ...Option) *Service {
	s := &Service{nftRepo: nftRepo, buyNFTFunc: buyNFTFunc}

	for _, fn := range opt {
		fn(s)
	}

	return s
}

func (s *Service) CreateNFT(ctx context.Context, userID uuid.UUID, nft *NFT) (string, error) {
	item, err := s.nftRepo.AddNFTItem(ctx, repository.AddNFTItemParams{
		Name:        nft.Name,
		Description: sql.NullString{String: nft.Description, Valid: len(nft.Description) > 0},
		Cover:       nft.ImageLink,
		Supply:      1, // TODO: multiple supplies are not supported yet
		BuyNowPrice: nft.BuyNowPrice,
		TokenURI:    nft.TokenURI,
	})
	if err != nil {
		return "", err
	}

	return item.ID.String(), nil
}

func (s *Service) BuyNFT(ctx context.Context, userID uuid.UUID, nftID uuid.UUID) error {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nftID)
	if err != nil {
		return fmt.Errorf("could not find NFT with id=%s: %w", nftID, err)
	}
	if item.OwnerID.Valid {
		return ErrAlreadySold
	}

	if err := s.buyNFTFunc(ctx, userID, item.BuyNowPrice, fmt.Sprintf("NFT purchase: %s", nftID)); err != nil {
		return fmt.Errorf("NFT purchase error: %w", err)
	}

	if err := s.nftRepo.UpdateNFTItemOwner(ctx, repository.UpdateNFTItemOwnerParams{
		ID:      nftID,
		OwnerID: uuid.NullUUID{UUID: userID, Valid: true},
	}); err != nil {
		// TODO: implement refund function or wrap operation into db transaction
		return fmt.Errorf("could not change NFT owner: %w", err)
	}

	return nil
}

func (s *Service) GetNFTs(ctx context.Context, limit, offset int32) ([]*NFT, error) {
	nftList, err := s.nftRepo.GetNFTItemsList(ctx, repository.GetNFTItemsListParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return castNFTRawListToNFTList(nftList), nil
}

func (s *Service) GetNFTsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, categoryID, limit, offset)
}

func (s *Service) GetNFTsByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, showID, limit, offset)
}

func (s *Service) GetNFTsByEpisodeID(ctx context.Context, episodeID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, episodeID, limit, offset)
}

func (s *Service) GetNFTsByRelationID(ctx context.Context, relID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	nftList, err := s.nftRepo.GetNFTItemsListByRelationID(ctx, repository.GetNFTItemsListByRelationIDParams{
		RelationID: relID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return castNFTRawListToNFTList(nftList), nil
}

func (s *Service) GetNFTsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	nftList, err := s.nftRepo.GetNFTItemsListByOwnerID(ctx, repository.GetNFTItemsListByOwnerIDParams{
		OwnerID: uuid.NullUUID{UUID: userID, Valid: true},
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return castNFTRawListToNFTList(nftList), nil
}

func (s *Service) GetNFTByID(ctx context.Context, nftID uuid.UUID) (*NFT, error) {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nftID)
	if err != nil {
		return nil, fmt.Errorf("could not find NFT with id=%s: %w", nftID, err)
	}

	return castNFTRawToNFT(item), nil
}

func (s *Service) GetCategories(ctx context.Context) ([]*Category, error) {
	clist, err := s.nftRepo.GetNFTCategoriesList(ctx)
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("could not get NFT categories list: %w", err)
	}

	return castCategoriesRawToCategories(clist), nil
}

func (s *Service) GetMainScreenCategory(ctx context.Context) (*Category, error) {
	c, err := s.nftRepo.GetMainNFTCategory(ctx)
	if err != nil {
		if db.IsNotFoundError(err) {
			clist, err := s.nftRepo.GetNFTCategoriesList(ctx)
			if err != nil {
				return nil, fmt.Errorf("could not get category to show on home screen: %w", err)
			}
			if len(clist) > 0 {
				return castCategoryRawToCategory(clist[rand.Int63n(int64(len(clist)))]), nil
			}

			return nil, nil
		}

		return nil, fmt.Errorf("could not get NFT categories list: %w", err)
	}

	return castCategoryRawToCategory(c), nil
}

func castCategoriesRawToCategories(clist []repository.NFTCategory) []*Category {
	res := make([]*Category, 0, len(clist))
	for _, i := range clist {
		res = append(res, castCategoryRawToCategory(i))
	}

	return res
}

func castCategoryRawToCategory(source repository.NFTCategory) *Category {
	return &Category{
		ID:    source.ID,
		Title: source.Title,
	}
}

func castNFTRawListToNFTList(source []repository.NFTItem) []*NFT {
	res := make([]*NFT, 0, len(source))
	for _, i := range source {
		res = append(res, castNFTRawToNFT(i))
	}

	return res
}

func castNFTRawToNFT(source repository.NFTItem) *NFT {
	nft := &NFT{
		ID:          source.ID,
		ImageLink:   source.Cover,
		Name:        source.Name,
		Description: source.Description.String,
		Supply:      int(source.Supply),
		BuyNowPrice: source.BuyNowPrice,
		TokenURI:    source.TokenURI,
	}

	return nft
}
