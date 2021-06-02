package client

import (
	"context"
	"encoding/json"

	"github.com/SatorNetwork/sator-api/svc/questions"
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
func (c *Client) GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]questions.Question, error) {
	res, err := c.s.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	list := make([]questions.Question, 0)
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// CheckAnswer ...
func (c *Client) CheckAnswer(ctx context.Context, aid uuid.UUID) (bool, error) {
	return c.s.CheckAnswer(ctx, aid)
}
