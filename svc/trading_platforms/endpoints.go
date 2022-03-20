package trading_platforms

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/utils"
	"github.com/SatorNetwork/sator-api/lib/validator"
)

type (
	Endpoints struct {
		CreateLink endpoint.Endpoint
		UpdateLink endpoint.Endpoint
		DeleteLink endpoint.Endpoint
		GetLinks   endpoint.Endpoint
	}

	service interface {
		CreateLink(ctx context.Context, req *CreateLinkRequest) (*Link, error)
		UpdateLink(ctx context.Context, req *UpdateLinkRequest) (*Link, error)
		DeleteLink(ctx context.Context, id uuid.UUID) error
		GetLinks(ctx context.Context, req *utils.PaginationRequest) ([]*Link, error)
	}

	Empty struct{}

	CreateLinkRequest struct {
		Title string `json:"title" validate:"required"`
		Link  string `json:"link" validate:"required"`
		Logo  string `json:"logo" validate:"required"`
	}

	UpdateLinkRequest struct {
		ID    uuid.UUID `json:"-" validate:"-"`
		Title string    `json:"title" validate:"required"`
		Link  string    `json:"link" validate:"required"`
		Logo  string    `json:"logo" validate:"required"`
	}

	DeleteLinkRequest struct {
		ID uuid.UUID `json:"id"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		CreateLink: MakeCreateLinkEndpoint(s, validateFunc),
		UpdateLink: MakeUpdateLinkEndpoint(s, validateFunc),
		DeleteLink: MakeDeleteLinkEndpoint(s),
		GetLinks:   MakeGetLinksEndpoint(s),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.CreateLink = mdw(e.CreateLink)
		}
	}

	return e
}

func MakeCreateLinkEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*CreateLinkRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to create-link-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		resp, err := s.CreateLink(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeUpdateLinkEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*UpdateLinkRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to update-link-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		resp, err := s.UpdateLink(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeDeleteLinkEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*DeleteLinkRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to delete-link-request")
		}

		if err := s.DeleteLink(ctx, typedReq.ID); err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeGetLinksEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		typedReq, ok := req.(utils.PaginationRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to get-links-request")
		}

		resp, err := s.GetLinks(ctx, &typedReq)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
