package client

import (
	"context"

	repository2 "github.com/SatorNetwork/sator-api/svc/challenge/repository"
	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		GetChallenges(ctx context.Context, arg repository2.GetChallengesParams) ([]repository2.Challenge, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetListByShowID returns challenges list filtered by show id
func (c *Client) GetListByShowID(ctx context.Context, showID uuid.UUID, page, itemsPerPage int32) (interface{}, error) {
	limit := itemsPerPage
	offset := limit * (page - 1)
	return c.s.GetChallenges(ctx, repository2.GetChallengesParams{
		ShowID: showID,
		Limit:  limit,
		Offset: offset,
	})
}
