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
		GetUsernameByID(ctx context.Context, uid uuid.UUID) (string, error)
	}
)

// New auth service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetUsernameByID ...
func (c *Client) GetUsernameByID(ctx context.Context, id uuid.UUID) (string, error) {
	return c.s.GetUsernameByID(ctx, id)
}
