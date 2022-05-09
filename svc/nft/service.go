package nft

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/db"
	lib_solana "github.com/SatorNetwork/sator-api/lib/solana"
	"github.com/SatorNetwork/sator-api/svc/nft/repository"
)

type (
	// Service struct
	Service struct {
		nftRepo    nftRepository
		sc         solanaClient
		buyNFTFunc buyNFTFunction
	}

	NFT struct {
		ID          uuid.UUID
		OwnerID     *uuid.UUID
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
		Minted     int32

		BuyNowPrice float64

		AuctionParams *NFTAuctionParams
		// NFT payload, e.g.: link to the original file, etc
		TokenURI    string
		RelationIDs []uuid.UUID
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
		AddNFTItemOwner(ctx context.Context, arg repository.AddNFTItemOwnerParams) error
		GetNFTItemByID(ctx context.Context, nftItemID uuid.UUID) (repository.GetNFTItemByIDRow, error)
		GetNFTItemsList(ctx context.Context, arg repository.GetNFTItemsListParams) ([]repository.GetNFTItemsListRow, error)
		GetAllNFTItems(ctx context.Context, arg repository.GetAllNFTItemsParams) ([]repository.GetAllNFTItemsRow, error)
		GetNFTItemsListByRelationID(ctx context.Context, arg repository.GetNFTItemsListByRelationIDParams) ([]repository.GetNFTItemsListByRelationIDRow, error)
		GetNFTItemsListByOwnerID(ctx context.Context, arg repository.GetNFTItemsListByOwnerIDParams) ([]repository.NFTItem, error)
		GetNFTCategoriesList(ctx context.Context) ([]repository.NFTCategory, error)
		GetMainNFTCategory(ctx context.Context) (repository.NFTCategory, error)
		DoesUserOwnNFT(ctx context.Context, arg repository.DoesUserOwnNFTParams) (bool, error)
		UpdateNFTItem(ctx context.Context, arg repository.UpdateNFTItemParams) error
		DeleteNFTItemByID(ctx context.Context, id uuid.UUID) error
		AddNFTRelation(ctx context.Context, arg repository.AddNFTRelationParams) error
		DoesRelationIDHasRelationNFT(ctx context.Context, relationID uuid.UUID) (bool, error)
	}

	solanaClient interface {
		GetNFTsByWalletAddress(ctx context.Context, walletAddr string) ([]*lib_solana.ArweaveNFTMetadata, error)
	}

	NFTItemRow struct {
		ID             uuid.UUID      `json:"id"`
		OwnerID        uuid.NullUUID  `json:"owner_id"`
		Name           string         `json:"name"`
		Description    sql.NullString `json:"description"`
		Cover          string         `json:"cover"`
		Supply         int64          `json:"supply"`
		BuyNowPrice    float64        `json:"buy_now_price"`
		TokenURI       string         `json:"token_uri"`
		UpdatedAt      sql.NullTime   `json:"updated_at"`
		CreatedAt      time.Time      `json:"created_at"`
		CreatorAddress sql.NullString `json:"creator_address"`
		CreatorShare   sql.NullInt32  `json:"creator_share"`
		Minted         sql.NullInt32  `json:"minted"`
	}

	// Simple function
	buyNFTFunction func(ctx context.Context, uid uuid.UUID, amount float64, info string, creatorAddr string, creatorShare int32) error
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(nftRepo nftRepository, sc solanaClient, buyNFTFunc buyNFTFunction, opt ...Option) *Service {
	s := &Service{
		nftRepo:    nftRepo,
		sc:         sc,
		buyNFTFunc: buyNFTFunc,
	}

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
		Supply:      int64(nft.Supply),
		BuyNowPrice: nft.BuyNowPrice,
		TokenURI:    nft.TokenURI,
	})
	if err != nil {
		return "", err
	}
	for i := 0; i < len(nft.RelationIDs); i++ {
		err := s.nftRepo.AddNFTRelation(ctx, repository.AddNFTRelationParams{
			NFTItemID:  item.ID,
			RelationID: nft.RelationIDs[i],
		})
		if err != nil {
			return "", err
		}
	}

	return item.ID.String(), nil
}

func (s *Service) BuyNFT(ctx context.Context, userID uuid.UUID, nftID uuid.UUID) error {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nftID)
	if err != nil {
		return fmt.Errorf("could not find NFT with id=%s: %w", nftID, err)
	}
	if item.Supply < int64(item.Minted) {
		return ErrAlreadySold
	}

	if yes, _ := s.nftRepo.DoesUserOwnNFT(ctx, repository.DoesUserOwnNFTParams{
		UserID:    userID,
		NFTItemID: nftID,
	}); yes {
		return ErrAlreadyBought
	}

	if err := s.buyNFTFunc(
		ctx,
		userID,
		item.BuyNowPrice,
		fmt.Sprintf("NFT purchase: %s", nftID),
		item.CreatorAddress.String,
		item.CreatorShare.Int32,
	); err != nil {
		return fmt.Errorf("NFT purchase error: %w", err)
	}

	//TODO: if owner db.NotFoundErr{AddItemOwner}
	if err := s.nftRepo.AddNFTItemOwner(ctx, repository.AddNFTItemOwnerParams{
		NFTItemID: nftID,
		UserID:    userID,
	}); err != nil {
		// TODO: implement refund function or wrap operation into db transaction
		return fmt.Errorf("could not change NFT owner: %w", err)
	}

	return nil
}

func (s *Service) GetNFTs(ctx context.Context, limit, offset int32, withMinted bool) ([]*NFT, error) {
	var ls []NFTItemRow

	if withMinted {
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

		ls = make([]NFTItemRow, 0, len(nftList))
		for _, v := range nftList {
			ls = append(ls, NFTItemRow(v))
		}
	} else {
		nftList, err := s.nftRepo.GetAllNFTItems(ctx, repository.GetAllNFTItemsParams{
			Limit:  limit,
			Offset: offset,
		})
		if err != nil {
			if db.IsNotFoundError(err) {
				return nil, nil
			}
			return nil, err
		}

		ls = make([]NFTItemRow, 0, len(nftList))
		for _, v := range nftList {
			ls = append(ls, NFTItemRow(v))
		}
	}

	return castNFTRawListToNFTList(ls), nil
}

func (s *Service) GetNFTsByCategory(ctx context.Context, uid, categoryID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, uid, categoryID, limit, offset)
}

func (s *Service) GetNFTsByShowID(ctx context.Context, uid, showID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, uid, showID, limit, offset)
}

func (s *Service) GetNFTsByEpisodeID(ctx context.Context, uid, episodeID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	return s.GetNFTsByRelationID(ctx, uid, episodeID, limit, offset)
}

func (s *Service) GetNFTsByRelationID(ctx context.Context, uid, relID uuid.UUID, limit, offset int32) ([]*NFT, error) {
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

	ls := make([]NFTItemRow, 0, len(nftList))
	for _, v := range nftList {
		ls = append(ls, NFTItemRow(v))
	}

	result := castNFTRawListToNFTList(ls)
	for k, item := range result {
		// TODO: needs refactoring! This is for backward compatibility with the app
		if yes, _ := s.nftRepo.DoesUserOwnNFT(ctx, repository.DoesUserOwnNFTParams{
			UserID:    uid,
			NFTItemID: item.ID,
		}); yes {
			result[k].OwnerID = &uid
		}
	}

	return result, nil
}

func (s *Service) GetNFTsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*NFT, error) {
	nftList, err := s.nftRepo.GetNFTItemsListByOwnerID(ctx, repository.GetNFTItemsListByOwnerIDParams{
		OwnerID: userID,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		if db.IsNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	result := castNFTListToNFTList(nftList)
	for k := range result {
		result[k].OwnerID = &userID
	}

	return result, nil
}

func (s *Service) GetNFTByID(ctx context.Context, nftID, userID uuid.UUID) (*NFT, error) {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nftID)
	if err != nil {
		return nil, fmt.Errorf("could not find NFT with id=%s: %w", nftID, err)
	}

	// TODO: needs refactoring! This is for backward compatibility with the app
	if yes, _ := s.nftRepo.DoesUserOwnNFT(ctx, repository.DoesUserOwnNFTParams{
		UserID:    userID,
		NFTItemID: nftID,
	}); yes {
		return castNFTRawToNFTRow(item, userID), nil
	}

	return castNFTRawToNFTRow(item), nil
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

func castNFTRawListToNFTList(source []NFTItemRow) []*NFT {
	res := make([]*NFT, 0, len(source))
	for _, i := range source {
		res = append(res, castNFTRawToNFT(i))
	}

	return res
}

func castNFTListToNFTList(source []repository.NFTItem) []*NFT {
	res := make([]*NFT, 0, len(source))
	for _, i := range source {
		res = append(res, castNFTItemToNFT(i))
	}

	return res
}

func castNFTRawToNFT(source NFTItemRow) *NFT {
	nft := &NFT{
		ID:          source.ID,
		ImageLink:   source.Cover,
		Name:        source.Name,
		Description: fmt.Sprintf("Rarity: #%d of %d. %s", source.Minted.Int32+1, source.Supply, source.Description.String),
		Supply:      int(source.Supply),
		BuyNowPrice: source.BuyNowPrice,
		TokenURI:    source.TokenURI,
		Minted:      source.Minted.Int32,
	}

	// if source.OwnerID.Valid && source.OwnerID.UUID != uuid.Nil {
	// 	nft.OwnerID = &source.OwnerID.UUID
	// }

	return nft
}

func castNFTItemToNFT(source repository.NFTItem) *NFT {
	nft := &NFT{
		ID:          source.ID,
		ImageLink:   source.Cover,
		Name:        source.Name,
		Description: source.Description.String,
		Supply:      int(source.Supply),
		BuyNowPrice: source.BuyNowPrice,
		TokenURI:    source.TokenURI,
	}

	// if source.OwnerID.Valid && source.OwnerID.UUID != uuid.Nil {
	// 	nft.OwnerID = &source.OwnerID.UUID
	// }

	return nft
}

func castNFTRawToNFTRow(source repository.GetNFTItemByIDRow, ownerID ...uuid.UUID) *NFT {
	nft := &NFT{
		ID:          source.ID,
		ImageLink:   source.Cover,
		Name:        source.Name,
		Description: fmt.Sprintf("Rarity: #%d of %d. %s", source.Minted+1, source.Supply, source.Description.String),
		Supply:      int(source.Supply),
		BuyNowPrice: source.BuyNowPrice,
		TokenURI:    source.TokenURI,
		Minted:      source.Minted,
	}

	if len(ownerID) > 0 && ownerID[0] != uuid.Nil {
		nft.OwnerID = &ownerID[0]
	}

	return nft
}

func (s *Service) DeleteNFTItemByID(ctx context.Context, nftID uuid.UUID) error {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nftID)
	if err != nil {
		return fmt.Errorf("could not find NFT with id=%s: %w", nftID, err)
	}
	if item.Minted > 0 {
		return ErrAlreadyMinted
	}

	err = s.nftRepo.DeleteNFTItemByID(ctx, nftID)
	if err != nil {
		return fmt.Errorf("could not delete NFT with id=%s: %w", nftID, err)
	}

	return nil
}

func (s *Service) UpdateNFTItem(ctx context.Context, nft *NFT) error {
	item, err := s.nftRepo.GetNFTItemByID(ctx, nft.ID)
	if err != nil {
		return fmt.Errorf("could not find NFT with id=%s: %w", nft.ID, err)
	}
	if item.Minted > 0 {
		return ErrAlreadyMinted
	}

	err = s.nftRepo.UpdateNFTItem(ctx, repository.UpdateNFTItemParams{
		ID:          nft.ID,
		Cover:       nft.ImageLink,
		Name:        nft.Name,
		Description: sql.NullString{String: nft.Description, Valid: len(nft.Description) > 0},
		Supply:      int64(nft.Supply),
		BuyNowPrice: nft.BuyNowPrice,
		TokenURI:    nft.TokenURI,
	})
	if err != nil {
		return fmt.Errorf("could not update NFT with id=%s: %w", nft.ID, err)
	}
	for i := 0; i < len(nft.RelationIDs); i++ {
		err := s.nftRepo.AddNFTRelation(ctx, repository.AddNFTRelationParams{
			NFTItemID:  item.ID,
			RelationID: nft.RelationIDs[i],
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) DoesRelationIDHasNFT(ctx context.Context, relationID uuid.UUID) (bool, error) {
	hasRelationID, err := s.nftRepo.DoesRelationIDHasRelationNFT(ctx, relationID)
	if err != nil {
		return false, fmt.Errorf("could not check is there related NFT with relation id=%s: %w", relationID, err)
	}

	return hasRelationID, nil
}

func (s *Service) GetNFTsByWalletAddress(ctx context.Context, req *GetNFTsByWalletAddressRequest) ([]*lib_solana.ArweaveNFTMetadata, error) {
	nfts, err := s.sc.GetNFTsByWalletAddress(ctx, req.WalletAddr)
	if err != nil {
		return nil, errors.Wrap(err, "can't get nfts from solana blockchain")
	}
	return nfts, nil
}
