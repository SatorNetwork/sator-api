package nft

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of NFT service
	Endpoints struct {
		CreateNFT         endpoint.Endpoint
		GetNFTs           endpoint.Endpoint
		GetNFTsByCategory endpoint.Endpoint
		GetNFTsByShowID   endpoint.Endpoint
		GetNFTsByUserID   endpoint.Endpoint
		GetNFTByID        endpoint.Endpoint
		BuyNFT            endpoint.Endpoint
		GetCategories     endpoint.Endpoint
	}

	service interface {
		CreateNFT(ctx context.Context, userUid uuid.UUID, nft *NFT) (string, error)
		GetNFTs(ctx context.Context) ([]*NFT, error)
		GetNFTsByCategory(ctx context.Context, category string) ([]*NFT, error)
		GetNFTsByShowID(ctx context.Context, showId, episodeId string) ([]*NFT, error)
		GetNFTsByUserID(ctx context.Context, userId string) ([]*NFT, error)
		GetNFTByID(ctx context.Context, nftId string) (*NFT, error)
		BuyNFT(ctx context.Context, userUid uuid.UUID, nftId string) error
		GetCategories(ctx context.Context) ([]*Category, error)
	}

	TransportNFT struct {
		ImageLink string `json:"image_link"`

		Name        string            `json:"name"`
		Description string            `json:"description"`
		Tags        map[string]string `json:"tags"`
		// Supply - the number of copies that can be minted.
		Supply int `json:"supply"`
		// Royalties are optional and allow user to earn a percentage on secondary sales
		Royalties  float64 `json:"royalties"` // TODO(evg): add validation?
		Blockchain string  `json:"blockchain"`
		SellType   string  `json:"sell_type"`

		AuctionParams *TransportNFTAuctionParams `json:"auction_params"`
	}

	TransportNFTAuctionParams struct {
		StartingBid    float64 `json:"starting_bid"`
		StartTimestamp string  `json:"start_timestamp"`
		EndTimestamp   string  `json:"end_timestamp"`
	}

	GetNFTsByCategoryRequest struct {
		Category string `json:"category"`
	}

	GetNFTsByShowIDRequest struct {
		ShowID    string `json:"show_id"`
		EpisodeID string `json:"episode_id"`
	}

	GetNFTsByUserIDRequest struct {
		UserID string `json:"user_id"`
	}

	TransportCategory struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}

	Empty struct{}
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
		ImageLink:   n.ImageLink,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
		Supply:      n.Supply,
		Royalties:   n.Royalties,
		Blockchain:  n.Blockchain,
		SellType:    n.SellType,
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
		ImageLink:   n.ImageLink,
		Name:        n.Name,
		Description: n.Description,
		Tags:        n.Tags,
		Supply:      n.Supply,
		Royalties:   n.Royalties,
		Blockchain:  n.Blockchain,
		SellType:    n.SellType,
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
	// validateFunc := validator.ValidateStruct()

	e := Endpoints{
		CreateNFT:         MakeCreateNFTEndpoint(s),
		GetNFTs:           MakeGetNFTsEndpoint(s),
		GetNFTsByCategory: MakeGetNFTsByCategoryEndpoint(s),
		GetNFTsByShowID:   MakeGetNFTsByShowIDEndpoint(s),
		GetNFTsByUserID:   MakeGetNFTsByUserIDEndpoint(s),
		GetNFTByID:        MakeGetNFTByIDEndpoint(s),
		BuyNFT:            MakeBuyNFTEndpoint(s),
		GetCategories:     MakeGetCategoriesEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.CreateNFT = mdw(e.CreateNFT)
			e.GetNFTs = mdw(e.GetNFTs)
			e.GetNFTsByCategory = mdw(e.GetNFTsByCategory)
			e.GetNFTsByShowID = mdw(e.GetNFTsByShowID)
			e.GetNFTsByUserID = mdw(e.GetNFTsByUserID)
			e.GetNFTByID = mdw(e.GetNFTByID)
			e.BuyNFT = mdw(e.BuyNFT)
			e.GetCategories = mdw(e.GetCategories)
		}
	}

	return e
}

func MakeCreateNFTEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		transportNFT, ok := request.(TransportNFT)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: TransportNFT, got: %T", request)
		}

		nftID, err := s.CreateNFT(ctx, uid, transportNFT.ToServiceNFT())
		if err != nil {
			return nil, err
		}

		return nftID, nil
	}
}

func MakeGetNFTsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		nfts, err := s.GetNFTs(ctx)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs: %v", err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByCategoryEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByCategoryRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByCategoryRequest, got: %T", request)
		}

		nfts, err := s.GetNFTsByCategory(ctx, req.Category)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by category: %v, %v", req.Category, err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByShowIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByShowIDRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByShowIDRequest, got: %T", request)
		}

		nfts, err := s.GetNFTsByShowID(ctx, req.ShowID, req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by show & episode ids: %v - %v, err: %v", req.ShowID, req.EpisodeID, err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTsByUserIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*GetNFTsByUserIDRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: GetNFTsByUserIDRequest, got: %T", request)
		}

		nfts, err := s.GetNFTsByUserID(ctx, req.UserID)
		if err != nil {
			return nil, fmt.Errorf("can't get NFTs by user's ID: %v", err)
		}

		return FromServiceNFTs(nfts), nil
	}
}

func MakeGetNFTByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		nftId, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: string, got: %T", request)
		}

		nft, err := s.GetNFTByID(ctx, nftId)
		if err != nil {
			return nil, fmt.Errorf("can't get nft by id: %v, err: %v", nftId, err)
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

		nftId, ok := request.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected request type, want: string, got: %T", request)
		}

		err = s.BuyNFT(ctx, userUid, nftId)
		if err != nil {
			return nil, fmt.Errorf("can't buy nft by id: %v, err: %v", nftId, err)
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
