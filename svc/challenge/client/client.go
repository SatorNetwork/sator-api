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
		GetChallengesByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetListByShowID returns challenges list filtered by show id
func (c *Client) GetListByShowID(ctx context.Context, showID uuid.UUID, limit, offset int32) (interface{}, error) {
	if limit < 1 {
		limit = 20
	}
	return c.s.GetChallengesByShowID(ctx, showID, limit, offset)
}
