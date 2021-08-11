package client

import (
	"context"
	"encoding/json"

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

		GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) (interface{}, error)
		GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*challenge.Question, error)
		CheckAnswer(ctx context.Context, aid, uid uuid.UUID) (bool, error)
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

// GetQuestionsByChallengeID returns questions list filtered by challenge id
func (c *Client) GetQuestionsByChallengeID(ctx context.Context, id uuid.UUID) ([]challenge.Question, error) {
	res, err := c.s.GetQuestionsByChallengeID(ctx, id)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	list := make([]challenge.Question, 0)
	if err := json.Unmarshal(b, &list); err != nil {
		return nil, err
	}

	return list, nil
}

// CheckAnswer ...
func (c *Client) CheckAnswer(ctx context.Context, aid, uid uuid.UUID) (bool, error) {
	return c.s.CheckAnswer(ctx, aid, uid)
}

// GetOneRandomQuestionByChallengeID ...
func (c *Client) GetOneRandomQuestionByChallengeID(ctx context.Context, id uuid.UUID, excludeIDs ...uuid.UUID) (*challenge.Question, error) {
	return c.s.GetOneRandomQuestionByChallengeID(ctx, id, excludeIDs...)
}
