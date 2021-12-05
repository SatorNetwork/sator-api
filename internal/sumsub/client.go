package sumsub

import (
	"context"

	"github.com/google/uuid"
)

type (
	Client struct {
		s service
	}

	service interface {
		GetSDKAccessTokenByApplicantID(ctx context.Context, applicantID string) (string, error)
		GetSDKAccessTokenByUserID(ctx context.Context, userID string) (string, error)
		GetByExternalUserID(ctx context.Context, externalUserID uuid.UUID) (*Response, error)
	}
)

// NewClient challenges service client implementation
func NewClient(s service) *Client {
	return &Client{s: s}
}

// GetSDKAccessTokenByApplicantID ...
func (c *Client) GetSDKAccessTokenByApplicantID(ctx context.Context, applicantID string) (string, error) {
	token, err := c.s.GetSDKAccessTokenByApplicantID(ctx, applicantID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetSDKAccessTokenByUserID ...
func (c *Client) GetSDKAccessTokenByUserID(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := c.s.GetSDKAccessTokenByUserID(ctx, userID.String())
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetByExternalUserID ...
func (c *Client) GetByExternalUserID(ctx context.Context, userID uuid.UUID) (*Response, error) {
	resp, err := c.s.GetByExternalUserID(ctx, userID)
	if err != nil {
		return &Response{}, err
	}

	return resp, nil
}
