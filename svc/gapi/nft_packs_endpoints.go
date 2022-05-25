package gapi

import (
	"context"

	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// NFTPacksEndpoints contains all of the endpoints for the NFTPacks service.
	NFTPacksEndpoints struct {
		AddNFTPackEndpoint        endpoint.Endpoint
		DeleteNFTPackEndpoint     endpoint.Endpoint
		GetNFTPackEndpoint        endpoint.Endpoint
		GetNFTPacksListEndpoint   endpoint.Endpoint
		SoftDeleteNFTPackEndpoint endpoint.Endpoint
		UpdateNFTPackEndpoint     endpoint.Endpoint
	}

	nftPacksService interface {
		AddNFTPack(ctx context.Context, name string, price float64, dropChances DropChances) (*NFTPackInfo, error)
		DeleteNFTPack(ctx context.Context, id uuid.UUID) error
		GetNFTPack(ctx context.Context, id uuid.UUID) (*NFTPackInfo, error)
		GetNFTPacksList(ctx context.Context) ([]NFTPackInfo, error)
		SoftDeleteNFTPack(ctx context.Context, id uuid.UUID) error
		UpdateNFTPack(ctx context.Context, id uuid.UUID, name string, price float64, dropChances DropChances) (*NFTPackInfo, error)
	}
)

// MakeNFTPacksEndpoints returns all of the endpoints for the NFTPacks service.
func MakeNFTPacksEndpoints(s nftPacksService, m ...endpoint.Middleware) NFTPacksEndpoints {
	validateFunc := validator.ValidateStruct()

	e := NFTPacksEndpoints{
		AddNFTPackEndpoint:        MakeAddNFTPackEndpoint(s, validateFunc),
		DeleteNFTPackEndpoint:     MakeDeleteNFTPackEndpoint(s),
		GetNFTPackEndpoint:        MakeGetNFTPackEndpoint(s),
		GetNFTPacksListEndpoint:   MakeGetNFTPacksListEndpoint(s),
		SoftDeleteNFTPackEndpoint: MakeSoftDeleteNFTPackEndpoint(s),
		UpdateNFTPackEndpoint:     MakeUpdateNFTPackEndpoint(s, validateFunc),
	}

	for _, mw := range m {
		e.AddNFTPackEndpoint = mw(e.AddNFTPackEndpoint)
		e.DeleteNFTPackEndpoint = mw(e.DeleteNFTPackEndpoint)
		e.GetNFTPackEndpoint = mw(e.GetNFTPackEndpoint)
		e.GetNFTPacksListEndpoint = mw(e.GetNFTPacksListEndpoint)
		e.SoftDeleteNFTPackEndpoint = mw(e.SoftDeleteNFTPackEndpoint)
		e.UpdateNFTPackEndpoint = mw(e.UpdateNFTPackEndpoint)
	}

	return e
}

// AddNFTPackRequest collects the request parameters for the AddNFTPack method.
type AddNFTPackRequest struct {
	Name        string      `json:"name,omitempty" validate:"required"`
	Price       float64     `json:"price,omitempty" validate:"required,gte=0"`
	DropChances DropChances `json:"drop_chances,omitempty"`
}

// MakeAddNFTPackEndpoint returns an endpoint that invokes AddNFTPack on the service.
func MakeAddNFTPackEndpoint(s nftPacksService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddNFTPackRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		info, err := s.AddNFTPack(ctx, req.Name, req.Price, req.DropChances)
		if err != nil {
			return nil, err
		}

		return info, nil
	}
}

// MakeDeleteNFTPackEndpoint returns an endpoint that invokes DeleteNFTPack on the service.
func MakeDeleteNFTPackEndpoint(s nftPacksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := s.DeleteNFTPack(ctx, uuid.MustParse(request.(string))); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// MakeGetNFTPackEndpoint returns an endpoint that invokes GetNFTPack on the service.
func MakeGetNFTPackEndpoint(s nftPacksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		info, err := s.GetNFTPack(ctx, uuid.MustParse(request.(string)))
		if err != nil {
			return nil, err
		}

		return info, nil
	}
}

// MakeGetNFTPacksListEndpoint returns an endpoint that invokes GetNFTPacksList on the service.
func MakeGetNFTPacksListEndpoint(s nftPacksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		list, err := s.GetNFTPacksList(ctx)
		if err != nil {
			return nil, err
		}

		return list, nil
	}
}

// MakeSoftDeleteNFTPackEndpoint returns an endpoint that invokes SoftDeleteNFTPack on the service.
func MakeSoftDeleteNFTPackEndpoint(s nftPacksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := s.SoftDeleteNFTPack(ctx, uuid.MustParse(request.(string))); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

// UpdateNFTPackRequest collects the request parameters for the UpdateNFTPack method.
type UpdateNFTPackRequest struct {
	ID          string      `json:"id,omitempty" validate:"required,uuid"`
	Name        string      `json:"name,omitempty" validate:"omitempty,min=1,max=32"`
	Price       float64     `json:"price,omitempty" validate:"omitempty,min=0"`
	DropChances DropChances `json:"drop_chances,omitempty"`
}

// MakeUpdateNFTPackEndpoint returns an endpoint that invokes UpdateNFTPack on the service.
func MakeUpdateNFTPackEndpoint(s nftPacksService, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateNFTPackRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		info, err := s.UpdateNFTPack(ctx, uuid.MustParse(req.ID), req.Name, req.Price, req.DropChances)
		if err != nil {
			return nil, err
		}

		return info, nil
	}
}
