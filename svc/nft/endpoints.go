package nft

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/validator"

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
	}

	service interface {
		CreateNFT(ctx context.Context, userUid uuid.UUID, nft *NFT) (string, error)
		GetNFTs(ctx context.Context, limit, offset int32) ([]*NFT, error)
		GetNFTsByCategory(ctx context.Context, categoryID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByEpisodeID(ctx context.Context, episodeID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTsByUserID(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*NFT, error)
		GetNFTByID(ctx context.Context, nftID uuid.UUID) (*NFT, error)
		BuyNFT(ctx context.Context, userUid uuid.UUID, nftID uuid.UUID) error
		GetCategories(ctx context.Context) ([]*Category, error)
		GetMainScreenCategory(ctx context.Context) (*Category, error)
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}

	TransportNFT struct {
		ID          uuid.UUID         `json:"id"`
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
	}

	TransportNFTAuctionParams struct {
		StartingBid    float64 `json:"starting_bid"`
		StartTimestamp string  `json:"start_timestamp"`
		EndTimestamp   string  `json:"end_timestamp"`
	}

	GetNFTsByCategoryRequest struct {
		Category string `json:"category" validate:"required,uuid"`

		PaginationRequest
	}

	GetNFTsByShowIDRequest struct {
		ShowID string `json:"show_id" validate:"required,uuid"`

		PaginationRequest
	}

	GetNFTsByEpisodeIDRequest struct {
		EpisodeID string `json:"episode_id" validate:"required,uuid"`

		PaginationRequest
	}

	GetNFTsByUserIDRequest struct {
		UserID string `json:"user_id" validate:"required,uuid"`

		PaginationRequest
	}

	TransportCategory struct {
		ID    string          `json:"id"`
		Title string          `json:"title"`
		Items []*TransportNFT `json:"items,omitempty"`
	}

	Empty struct{}
)

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}
	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}
	return 0
}

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
		ImageLink:   n.ImageLink,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
		Supply:      n.Supply,
		Royalties:   n.Royalties,
		Blockchain:  n.Blockchain,
		SellType:    n.SellType,
		BuyNowPrice: n.BuyNowPrice,
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
		req, ok := request.(PaginationRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: PaginationRequest, got: %T", request)
		}
		if err := v(req); err != nil {
			return nil, err
		}

		nfts, err := s.GetNFTs(ctx, req.Limit(), req.Offset())
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

		nfts, err := s.GetNFTsByCategory(ctx, uuid.MustParse(req.Category), req.Limit(), req.Offset())
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

		nfts, err := s.GetNFTsByShowID(ctx, uuid.MustParse(req.ShowID), req.Limit(), req.Offset())
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

		nfts, err := s.GetNFTsByEpisodeID(ctx, uuid.MustParse(req.EpisodeID), req.Limit(), req.Offset())
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

		nft, err := s.GetNFTByID(ctx, uuid.MustParse(nftID))
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
			return nil, fmt.Errorf("can't buy nft by id: %v, err: %v", nftID, err)
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

		nfts, err := s.GetNFTsByCategory(ctx, cat.ID, 3, 0)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by category: %v, %v", cat.ID.String(), err)
		}

		category := FromServiceCategory(cat)
		category.Items = FromServiceNFTs(nfts)

		return category, nil
	}
}
