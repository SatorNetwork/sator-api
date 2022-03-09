package nft

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/SatorNetwork/sator-api/lib/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of NFT service
	Endpoints struct {
		CreateNFT          endpoint.Endpoint
		GetNFTs            endpoint.Endpoint
		GetNFTsByCategory  endpoint.Endpoint
		GetNFTsByShowID    endpoint.Endpoint
		GetNFTsByEpisodeID endpoint.Endpoint
		GetNFTsByUserID    endpoint.Endpoint
		GetNFTByID         endpoint.Endpoint
		BuyNFT             endpoint.Endpoint
		GetCategories      endpoint.Endpoint
		GetMainScreenData  endpoint.Endpoint
		DeleteNFTItemByID  endpoint.Endpoint
		UpdateNFTItem      endpoint.Endpoint
	}

	service interface {
		CreateNFT(ctx context.Context, userUid uuid.UUID, nft *NFT) (string, error)
		GetNFTs(ctx context.Context, limit, offset int32, withMinted bool) ([]*NFT, error)
		GetNFTsByCategory(ctx context.Context, uid, categoryID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByShowID(ctx context.Context, uid, showID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByEpisodeID(ctx context.Context, uid, episodeID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTByID(ctx context.Context, nftID, userID uuid.UUID) (*NFT, error)
		BuyNFT(ctx context.Context, userUid uuid.UUID, nftID uuid.UUID) error
		GetCategories(ctx context.Context) ([]*Category, error)
		GetMainScreenCategory(ctx context.Context) (*Category, error)
		DeleteNFTItemByID(ctx context.Context, nftID uuid.UUID) error
		UpdateNFTItem(ctx context.Context, nft *NFT) error
		GetNFTsByRelationID(ctx context.Context, uid, relID uuid.UUID, limit, offset int32) ([]*NFT, error)
	}

	TransportNFT struct {
		ID          uuid.UUID         `json:"id"`
		OwnerID     *uuid.UUID        `json:"owner_id,omitempty"`
		ImageLink   string            `json:"image_link"`
		Name        string            `json:"name" validate:"required"`
		Description string            `json:"description"`
		Tags        map[string]string `json:"tags,omitempty"`
		// Supply - the number of copies that can be minted.
		Supply int `json:"supply"`
		// Royalties are optional and allow user to earn a percentage on secondary sales
		Royalties   float64 `json:"royalties"` // TODO(evg): add validation?
		Blockchain  string  `json:"blockchain"`
		SellType    string  `json:"sell_type"`
		BuyNowPrice float64 `json:"buy_now_price"`
		TokenURI    string  `json:"token_uri" validate:"required"`

		AuctionParams *TransportNFTAuctionParams `json:"auction_params"`
		RelationIDs   []uuid.UUID                `json:"relation_ids"`
	}

	TransportNFTAuctionParams struct {
		StartingBid    float64 `json:"starting_bid"`
		StartTimestamp string  `json:"start_timestamp"`
		EndTimestamp   string  `json:"end_timestamp"`
	}

	GetNFTsByCategoryRequest struct {
		Category string `json:"category" validate:"required,uuid"`

		utils.PaginationRequest
	}

	GetNFTsWithFilterRequest struct {
		RelationID string `json:"relation_id,omitempty"`

		utils.PaginationRequest
	}

	GetNFTsByShowIDRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`

		utils.PaginationRequest
	}

	GetNFTsByEpisodeIDRequest struct {
		EpisodeID string `json:"episode_id" validate:"required,uuid"`

		utils.PaginationRequest
	}

	GetNFTsByUserIDRequest struct {
		UserID string `json:"user_id" validate:"required,uuid"`

		utils.PaginationRequest
	}

	TransportCategory struct {
		ID    string          `json:"id"`
		Title string          `json:"title"`
		Items []*TransportNFT `json:"items,omitempty"`
	}

	Empty struct{}

	UpdateNFTRequest struct {
		ID          uuid.UUID `json:"id"`
		ImageLink   string    `json:"image_link"`
		Name        string    `json:"name" validate:"required"`
		Description string    `json:"description"`
		Supply      int       `json:"supply"`
		BuyNowPrice float64   `json:"buy_now_price"`
		TokenURI    string    `json:"token_uri" validate:"required"`
	}
)

func FromServiceNFTs(nfts []*NFT) []*TransportNFT {
	transportNFTs := make([]*TransportNFT, 0, len(nfts))
	for _, n := range nfts {
		transportNFTs = append(transportNFTs, FromServiceNFT(n))
	}

	return transportNFTs
}

func FromServiceNFT(n *NFT) *TransportNFT {
	nft := &TransportNFT{
		ID:          n.ID,
		OwnerID:     n.OwnerID,
		ImageLink:   n.ImageLink,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
		Supply:      n.Supply,
		Royalties:   n.Royalties,
		Blockchain:  n.Blockchain,
		SellType:    n.SellType,
		BuyNowPrice: n.BuyNowPrice,
		TokenURI:    n.TokenURI,
	}
	if n.AuctionParams != nil {
		nft.AuctionParams = FromServiceNFTAuctionParams(n.AuctionParams)
	}

	return nft
}

func FromServiceNFTAuctionParams(a *NFTAuctionParams) *TransportNFTAuctionParams {
	return &TransportNFTAuctionParams{
		StartingBid:    a.StartingBid,
		StartTimestamp: a.StartTimestamp,
		EndTimestamp:   a.EndTimestamp,
	}
}

func (n *TransportNFT) ToServiceNFT() *NFT {
	nft := &NFT{
		ID:          n.ID,
		ImageLink:   n.ImageLink,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
		Supply:      n.Supply,
		Royalties:   n.Royalties,
		Blockchain:  n.Blockchain,
		SellType:    n.SellType,
		BuyNowPrice: n.BuyNowPrice,
		TokenURI:    n.TokenURI,
		RelationIDs: n.RelationIDs,
	}
	if n.AuctionParams != nil {
		nft.AuctionParams = n.AuctionParams.ToServiceNFTAuctionParams()
	}

	return nft
}

func (a *TransportNFTAuctionParams) ToServiceNFTAuctionParams() *NFTAuctionParams {
	return &NFTAuctionParams{
		StartingBid:    a.StartingBid,
		StartTimestamp: a.StartTimestamp,
		EndTimestamp:   a.EndTimestamp,
	}
}

func FromServiceCategory(c *Category) *TransportCategory {
	return &TransportCategory{
		ID:    c.ID.String(),
		Title: c.Title,
	}
}

func FromServiceCategories(categories []*Category) []*TransportCategory {
	transportCategories := make([]*TransportCategory, 0, len(categories))
	for _, c := range categories {
		transportCategories = append(transportCategories, FromServiceCategory(c))
	}

	return transportCategories
}

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		CreateNFT:          MakeCreateNFTEndpoint(s, validateFunc),
		GetNFTs:            MakeGetNFTsEndpoint(s, validateFunc),
		GetNFTsByCategory:  MakeGetNFTsByCategoryEndpoint(s, validateFunc),
		GetNFTsByShowID:    MakeGetNFTsByShowIDEndpoint(s, validateFunc),
		GetNFTsByEpisodeID: MakeGetNFTsByEpisodeIDEndpoint(s, validateFunc),
		GetNFTsByUserID:    MakeGetNFTsByUserIDEndpoint(s, validateFunc),
		GetNFTByID:         MakeGetNFTByIDEndpoint(s),
		BuyNFT:             MakeBuyNFTEndpoint(s),
		GetCategories:      MakeGetCategoriesEndpoint(s),
		GetMainScreenData:  MakeGetMainScreenDataEndpoint(s),
		DeleteNFTItemByID:  MakeDeleteNFTItemByIDEndpoint(s),
		UpdateNFTItem:      MakeUpdateNFTItemEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.CreateNFT = mdw(e.CreateNFT)
			e.GetNFTs = mdw(e.GetNFTs)
			e.GetNFTsByCategory = mdw(e.GetNFTsByCategory)
			e.GetNFTsByShowID = mdw(e.GetNFTsByShowID)
			e.GetNFTsByEpisodeID = mdw(e.GetNFTsByEpisodeID)
			e.GetNFTsByUserID = mdw(e.GetNFTsByUserID)
			e.GetNFTByID = mdw(e.GetNFTByID)
			e.BuyNFT = mdw(e.BuyNFT)
			e.GetCategories = mdw(e.GetCategories)
			e.GetMainScreenData = mdw(e.GetMainScreenData)
			e.DeleteNFTItemByID = mdw(e.DeleteNFTItemByID)
			e.UpdateNFTItem = mdw(e.UpdateNFTItem)
		}
	}

	return e
}

func MakeCreateNFTEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		transportNFT, ok := request.(TransportNFT)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: TransportNFT, got: %T", request)
		}
		if err := v(transportNFT); err != nil {
			return nil, err
		}

		nftID, err := s.CreateNFT(ctx, uid, transportNFT.ToServiceNFT())
		if err != nil {
			return nil, err
		}

		return nftID, nil
	}
}

func MakeGetNFTsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(*GetNFTsWithFilterRequest)
		if err := v(req); err != nil {
			return nil, err
		}
		if req.RelationID != "" {
			uid, err := jwt.UserIDFromContext(ctx)
			if err != nil {
				return nil, fmt.Errorf("could not get user profile id: %w", err)
			}

			relationID, err := uuid.Parse(req.RelationID)
			if err != nil {
				return nil, fmt.Errorf("could not get relation id: %w", err)
			}

			nfts, err := s.GetNFTsByRelationID(ctx, uid, relationID, req.Limit(), req.Offset())
			if err != nil {
				return nil, err
			}

			return FromServiceNFTs(nfts), nil
		}

		nfts, err := s.GetNFTs(ctx, req.Limit(), req.Offset(), rbac.IsCurrentUserHasRole(ctx, rbac.RoleAdmin, rbac.RoleContentManager))
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs: %v", err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByCategoryEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByCategoryRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByCategoryRequest, got: %T", request)
		}
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nfts, err := s.GetNFTsByCategory(ctx, uid, uuid.MustParse(req.Category), req.Limit(), req.Offset())
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by category: %v, %v", req.Category, err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByShowIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByShowIDRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByShowIDRequest, got: %T", request)
		}
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nfts, err := s.GetNFTsByShowID(ctx, uid, uuid.MustParse(req.ShowID), req.Limit(), req.Offset())
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by show id: %v, err: %v", req.ShowID, err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByEpisodeIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByEpisodeIDRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByEpisodeIDRequest, got: %T", request)
		}
		if err := v(req); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nfts, err := s.GetNFTsByEpisodeID(ctx, uid, uuid.MustParse(req.EpisodeID), req.Limit(), req.Offset())
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by episode id: %v, err: %v", req.EpisodeID, err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByUserIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByUserIDRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByUserIDRequest, got: %T", request)
		}
		if err := v(req); err != nil {
			return nil, err
		}

		nfts, err := s.GetNFTsByUserID(ctx, uuid.MustParse(req.UserID), req.Limit(), req.Offset())
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by user's ID: %v", err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		nftID, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: string, got: %T", request)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nft, err := s.GetNFTByID(ctx, uuid.MustParse(nftID), uid)
		if err != nil {
			return nil, fmt.Errorf("can't get nft by id: %v, err: %v", nftID, err)
		}

		return FromServiceNFT(nft), nil
	}
}

func MakeBuyNFTEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userUid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nftID, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: string, got: %T", request)
		}

		err = s.BuyNFT(ctx, userUid, uuid.MustParse(nftID))
		if err != nil {
			return nil, err
		}

		return Empty{}, nil
	}
}

func MakeGetCategoriesEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		categories, err := s.GetCategories(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't get categories: %v", err)
		}

		return FromServiceCategories(categories), nil
	}
}

func MakeGetMainScreenDataEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		cat, err := s.GetMainScreenCategory(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not found any category to show on main screen: %v", err)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		nfts, err := s.GetNFTsByCategory(ctx, uid, cat.ID, 3, 0)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by category: %v, %v", cat.ID.String(), err)
		}

		category := FromServiceCategory(cat)
		category.Items = FromServiceNFTs(nfts)

		return category, nil
	}
}

func MakeDeleteNFTItemByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		err = s.DeleteNFTItemByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeUpdateNFTItemEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(UpdateNFTRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		err := s.UpdateNFTItem(ctx, &NFT{
			ID:          req.ID,
			ImageLink:   req.ImageLink,
			Name:        req.Name,
			Description: req.Description,
			Supply:      req.Supply,
			BuyNowPrice: req.BuyNowPrice,
			TokenURI:    req.TokenURI,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
