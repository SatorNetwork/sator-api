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
		CheckAnswer(ctx context.Context, id uuid.UUID) (bool, error)
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

// CheckAnswer ...
func (c *Client) CheckAnswer(ctx context.Context, aid uuid.UUID) (interface{}, error) {
	return c.s.CheckAnswer(ctx, aid)
}
