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
		AddReward(ctx context.Context, uid uuid.UUID, amount float64, qid uuid.UUID) error
		GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// AddReward ...
func (c *Client) AddReward(ctx context.Context, userID uuid.UUID, amount float64, quizID uuid.UUID) error {
	return c.s.AddReward(ctx, userID, amount, quizID)
}

// GetUserRewards ...
func (c *Client) GetUserRewards(ctx context.Context, userID uuid.UUID) (float64, error) {
	return c.s.GetUserRewards(ctx, userID)
}
