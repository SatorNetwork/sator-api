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
		GetByID(ctx context.Context, id uuid.UUID) (challenge.Challenge, error)
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

// GetChallengeByID returns Challenge struct
func (c *Client) GetChallengeByID(ctx context.Context, challengeID uuid.UUID) (challenge.Challenge, error) {
	return c.s.GetByID(ctx, challengeID)
}
