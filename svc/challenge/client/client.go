package client

import (
	"context"

	"github.com/SatorNetwork/sator-api/svc/challenge"
	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		GetListByShowID(ctx context.Context, showID uuid.UUID, page, itemsPerPage int64) ([]challenge.Challenge, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetListByShowID returns challenges list filtered by show id
func (c *Client) GetListByShowID(ctx context.Context, showID uuid.UUID, page, itemsPerPage int64) (interface{}, error) {
	limit := itemsPerPage
	offset := limit * (page - 1)
	return c.s.GetListByShowID(ctx, showID, limit, offset)
}
