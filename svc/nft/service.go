package nft

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Service struct
	Service struct{}

	NFT struct {
		ImageLink string

		Name        string
		Description string
		Tags        map[string]string
		// Supply - the number of copies that can be minted.
		Supply int
		// Royalties are optional and allow user to earn a percentage on secondary sales
		Royalties  float64 // TODO(evg): add validation?
		Blockchain string  // TODO(evg): replace with enum?
		SellType   string  // TODO(evg): replace with enum?

		AuctionParams *NFTAuctionParams
	}

	NFTAuctionParams struct {
		StartingBid    float64
		StartTimestamp string // TODO(evg): replace with time.Time?
		EndTimestamp   string // TODO(evg): replace with time.Time?
	}

	// Option func to set custom service options
	Option func(*Service)
)

var (
	rfc3339Timestamp = "2006-01-02T15:04:05Z07:00"

	fakeNFT = NFT{
		ImageLink:   "https://sator-dev-storage.nyc3.cdn.digitaloceanspaces.com/uploads/6e3500c8-df21-4279-a092-33c7a0d73e90.png",
		Name:        "test name",
		Description: "test description",
		Tags: map[string]string{
			"test key": "test val",
		},
		Supply:     1,
		Royalties:  2,
		Blockchain: "Ethereum",
		SellType:   "Auction",
		AuctionParams: &NFTAuctionParams{
			StartingBid:    1,
			StartTimestamp: rfc3339Timestamp,
			EndTimestamp:   rfc3339Timestamp,
		},
	}
)

// NewService is a factory function,
// returns a new instance of the Service interface implementation
func NewService(opt ...Option) *Service {
	s := &Service{}

	for _, fn := range opt {
		fn(s)
	}

	return s
}

func (s *Service) CreateNFT(ctx context.Context, userUid uuid.UUID, nft *NFT) (string, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return uuid.New().String(), nil
}

func (s *Service) GetNFTs(ctx context.Context) ([]*NFT, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return []*NFT{&fakeNFT, &fakeNFT, &fakeNFT}, nil
}

func (s *Service) GetNFTsByCategory(ctx context.Context, category string) ([]*NFT, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return []*NFT{&fakeNFT, &fakeNFT, &fakeNFT}, nil
}

func (s *Service) GetNFTsByShowID(ctx context.Context, showId, episodeId string) ([]*NFT, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return []*NFT{&fakeNFT, &fakeNFT, &fakeNFT}, nil
}

func (s *Service) GetNFTsByUserID(ctx context.Context, userId string) ([]*NFT, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return []*NFT{&fakeNFT, &fakeNFT, &fakeNFT}, nil
}

func (s *Service) GetNFTByID(ctx context.Context, nftId string) (*NFT, error) {
	// TODO(evg): implement when NFT SDK will be ready
	return &fakeNFT, nil
}

func (s *Service) BuyNFT(ctx context.Context, userUid uuid.UUID, nftId string) error {
	// TODO(evg): implement when NFT SDK will be ready
	return nil
}
