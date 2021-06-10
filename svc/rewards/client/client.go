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
		AddTransaction(ctx context.Context, uid uuid.UUID, amount float64, qid uuid.UUID, trType int32) error
		GetUserRewards(ctx context.Context, uid uuid.UUID) (float64, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// AddTransaction ...
func (c *Client) AddTransaction(ctx context.Context, userID uuid.UUID, amount float64, quizID uuid.UUID, trType int32) error {
	return c.s.AddTransaction(ctx, userID, amount, quizID, trType)
}

// GetUserRewards ...
func (c *Client) GetUserRewards(ctx context.Context, userID uuid.UUID) (float64, error) {
	return c.s.GetUserRewards(ctx, userID)
}
