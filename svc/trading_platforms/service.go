package trading_platforms

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/utils"
	trading_platforms_repository "github.com/SatorNetwork/sator-api/svc/trading_platforms/repository"
)

type (
	Service struct {
		tpr tradingPlatformRepository
	}

	tradingPlatformRepository interface {
		CreateTradingPlatformLink(ctx context.Context, arg trading_platforms_repository.CreateTradingPlatformLinkParams) (trading_platforms_repository.TradingPlatformLink, error)
		UpdateTradingPlatformLink(ctx context.Context, arg trading_platforms_repository.UpdateTradingPlatformLinkParams) (trading_platforms_repository.TradingPlatformLink, error)
		DeleteTradingPlatformLink(ctx context.Context, id uuid.UUID) error
		GetTradingPlatformLinks(ctx context.Context, arg trading_platforms_repository.GetTradingPlatformLinksParams) ([]trading_platforms_repository.TradingPlatformLink, error)
	}

	Link struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Link  string `json:"link"`
		Logo  string `json:"logo"`
	}
)

func NewService(tpr tradingPlatformRepository) *Service {
	s := &Service{
		tpr: tpr,
	}

	return s
}

func (s *Service) CreateLink(ctx context.Context, req *CreateLinkRequest) (*Link, error) {
	resp, err := s.tpr.CreateTradingPlatformLink(ctx, trading_platforms_repository.CreateTradingPlatformLinkParams{
		Title: req.Title,
		Link:  req.Link,
		Logo:  req.Logo,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't create trading platform link")
	}

	return NewLinkFromSQLC(&resp), nil
}

func (s *Service) UpdateLink(ctx context.Context, req *UpdateLinkRequest) (*Link, error) {
	resp, err := s.tpr.UpdateTradingPlatformLink(ctx, trading_platforms_repository.UpdateTradingPlatformLinkParams{
		ID:    req.ID,
		Title: req.Title,
		Link:  req.Link,
		Logo:  req.Logo,
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't update trading platform link")
	}

	return NewLinkFromSQLC(&resp), nil
}

func (s *Service) DeleteLink(ctx context.Context, id uuid.UUID) error {
	if err := s.tpr.DeleteTradingPlatformLink(ctx, id); err != nil {
		return errors.Wrap(err, "can't delete trading platform link")
	}

	return nil
}

func (s *Service) GetLinks(ctx context.Context, req *utils.PaginationRequest) ([]*Link, error) {
	links, err := s.tpr.GetTradingPlatformLinks(ctx, trading_platforms_repository.GetTradingPlatformLinksParams{
		Limit:  req.Limit(),
		Offset: req.Offset(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "can't get trading platform links")
	}

	return NewLinksFromSQLC(links), nil
}

func NewLinkFromSQLC(link *trading_platforms_repository.TradingPlatformLink) *Link {
	return &Link{
		ID:    link.ID.String(),
		Title: link.Title,
		Link:  link.Link,
		Logo:  link.Logo,
	}
}

func NewLinksFromSQLC(sqlcLinks []trading_platforms_repository.TradingPlatformLink) []*Link {
	links := make([]*Link, 0, len(sqlcLinks))
	for _, sqlcLink := range sqlcLinks {
		links = append(links, NewLinkFromSQLC(&sqlcLink))
	}

	return links
}
