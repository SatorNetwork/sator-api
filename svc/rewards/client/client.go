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
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// AddReward ...
func (c *Client) AddReward(ctx context.Context, userID uuid.UUID, amount float64, quizID uuid.UUID) error {
	return nil
}
