package invitations

import (
	"context"
	"fmt"
	"strings"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		SendInvitation endpoint.Endpoint
	}

	service interface {
		SendInvitation(ctx context.Context, invitedByID uuid.UUID, invitedByUsername, inviteeEmail string) error
	}

	// SendInvitationRequest struct
	SendInvitationRequest struct {
		InviteeEmail string `json:"email" validate:"required,email"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		SendInvitation: MakeSendInvitationEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.SendInvitation = mdw(e.SendInvitation)
		}
	}

	return e
}

// MakeSendInvitationEndpoint ...
func MakeSendInvitationEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		req := request.(SendInvitationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		// normalize email address
		req.InviteeEmail = strings.ToLower(req.InviteeEmail)

		// invited by current user
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}
		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		if err := s.SendInvitation(ctx, uid, username, req.InviteeEmail); err != nil {
			return nil, err
		}

		return true, nil
	}
}
