package client

import (
	"context"

	"github.com/SatorNetwork/sator-api/svc/invitations"

	"github.com/google/uuid"
)

type (
	// Client struct
	Client struct {
		s service
	}

	service interface {
		AcceptInvitation(ctx context.Context, inviteeID uuid.UUID, inviteeEmail string) error
		GetInvitations(ctx context.Context) ([]invitations.Invitation, error)
		IsEmailInvited(ctx context.Context, inviteeEmail string) (bool, error)
	}
)

// New challenges service client implementation
func New(s service) *Client {
	return &Client{s: s}
}

// GetInvitations returns list invitations.
func (c *Client) GetInvitations(ctx context.Context) ([]invitations.Invitation, error) {
	resp, err := c.s.GetInvitations(ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// AcceptInvitation used to accept invitation and store invitee ID and email.
func (c *Client) AcceptInvitation(ctx context.Context, inviteeID uuid.UUID, inviteeEmail string) error {
	err := c.s.AcceptInvitation(ctx, inviteeID, inviteeEmail)
	if err != nil {
		return err
	}

	return nil
}

// IsEmailInvited returns true if email invited, false if not.
func (c *Client) IsEmailInvited(ctx context.Context, email string) (bool, error) {
	resp, err := c.s.IsEmailInvited(ctx, email)
	if err != nil {
		return false, err
	}

	return resp, nil
}
