package client

import (
	"context"

	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		DoesRelationIDHasNFT(ctx context.Context, relationID uuid.UUID) (bool, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// DoesRelationIDHasNFT ...
func (c *Client) DoesRelationIDHasNFT(ctx context.Context, relationID uuid.UUID) (bool, error) {
	return c.s.DoesRelationIDHasNFT(ctx, relationID)
}
