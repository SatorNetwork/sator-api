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
		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetQuestionsByChallengeID returns questions list filtered by challenge id
func (c *Client) GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error) {
	return c.s.GetQuestionsByChallengeID(ctx, id)
}
