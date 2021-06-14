package client

import (
	"context"

	"github.com/SatorNetwork/sator-api/svc/rewards"
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

// AddDepositTransaction ...
func (c *Client) AddDepositTransaction(ctx context.Context, userID uuid.UUID, amount float64, quizID uuid.UUID) error {
	return c.s.AddTransaction(ctx, userID, amount, quizID, rewards.TransactionTypeDeposit)
}

// AddWithdrawTransaction ...
func (c *Client) AddWithdrawTransaction(ctx context.Context, userID uuid.UUID, amount float64) error {
	return c.s.AddTransaction(ctx, userID, amount, uuid.Nil, rewards.TransactionTypeWithdraw)
}

// GetUserRewards ...
func (c *Client) GetUserRewards(ctx context.Context, userID uuid.UUID) (float64, error) {
	return c.s.GetUserRewards(ctx, userID)
}
