package sumsub

import (
	"context"

	"github.com/google/uuid"
)

type (
	Client struct {
		s service
	}

	service interface {
		GetSDKAccessTokenByUserID(ctx context.Context, userID, levelName string) (string, error)
	}
)

// NewClient challenges service client implementation
func NewClient(s service) *Client {
	return &Client{s: s}
}

// GetSDKAccessTokenByUserID ...
func (c *Client) GetSDKAccessTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := c.s.GetSDKAccessTokenByUserID(ctx, userID.String(), BasicKYCLevel)
	if err != nil {
		return "", err
	}

	return token, nil
}
