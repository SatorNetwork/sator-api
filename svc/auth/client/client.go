package client

import (
	"context"
	"crypto/rsa"
	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		GetUsernameByID(ctx context.Context, uid uuid.UUID) (string, error)
		GetPublicKey(ctx context.Context, userID uuid.UUID) (*rsa.PublicKey, error)
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

func (c *Client) GetPublicKey(ctx context.Context, userID uuid.UUID) (*rsa.PublicKey, error) {
	return c.s.GetPublicKey(ctx, userID)
}
