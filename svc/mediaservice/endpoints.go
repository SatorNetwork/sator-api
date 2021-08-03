package mediaservice

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddItem        endpoint.Endpoint
		GetItemByID    endpoint.Endpoint
		GetItemsList   endpoint.Endpoint
		DeleteItemByID endpoint.Endpoint
	}

	service interface {
		AddItem(ctx context.Context, it Item, file io.ReadSeeker, fileHeader *multipart.FileHeader) (Item, error)
		GetItemByID(ctx context.Context, id uuid.UUID) (Item, error)
		GetItemsList(ctx context.Context, limit, offset int32) ([]Item, error)
		DeleteItemByID(ctx context.Context, id uuid.UUID) error
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}

	// AddItemRequest struct
	AddItemRequest struct {
		File       io.ReadSeeker
		FileHeader *multipart.FileHeader
	}
)

// Limit of items
func (r PaginationRequest) Limit() int32 {
	if r.ItemsPerPage > 0 {
		return r.ItemsPerPage
	}

	return 20
}

// Offset items
func (r PaginationRequest) Offset() int32 {
	if r.Page > 1 {
		return (r.Page - 1) * r.Limit()
	}

	return 0
}

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		AddItem:        MakeAddItemEndpoint(s, validateFunc),
		GetItemByID:    MakeGetItemByIDEndpoint(s),
		GetItemsList:   MakeGetItemsListEndpoint(s, validateFunc),
		DeleteItemByID: MakeDeleteItemByIDEndpoint(s),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddItem = mdw(e.AddItem)
			e.GetItemByID = mdw(e.GetItemByID)
			e.GetItemsList = mdw(e.GetItemsList)
			e.DeleteItemByID = mdw(e.DeleteItemByID)
		}
	}

	return e
}

// MakeGetItemsListEndpoint ...
func MakeGetItemsListEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetItemsList(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeAddItemEndpoint ...
func MakeAddItemEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddItemRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddItem(ctx, Item{
			Filename: req.FileHeader.Filename,
		}, req.File, req.FileHeader)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetItemByIDEndpoint ...
func MakeGetItemByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("%w item id: %v", ErrInvalidParameter, err)
		}

		resp, err := s.GetItemByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteItemByIDEndpoint ...
func MakeDeleteItemByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get item id: %w", err)
		}

		err = s.DeleteItemByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
