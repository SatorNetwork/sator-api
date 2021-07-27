package invitations

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		SendInvitation endpoint.Endpoint
	}

	service interface {
		SendInvitation(ctx context.Context, invitedByID uuid.UUID, inviteeEmail string) error
	}

	// SendInvitationRequest struct
	SendInvitationRequest struct {
		InvitedByID  string `json:"invited_by_id" validate:"required,uuid"`
		InviteeEmail string `json:"invitee_email" validate:"required,gt=0"`
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
		req := request.(SendInvitationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		invitedByID, err := uuid.Parse(req.InvitedByID)
		if err != nil {
			return nil, fmt.Errorf("could not get invited by id: %w", err)
		}

		err = s.SendInvitation(ctx, invitedByID, req.InviteeEmail)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
