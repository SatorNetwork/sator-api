package qrcodes

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/jwt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of qrcode service
	Endpoints struct {
		GetDataByQRCodeID endpoint.Endpoint
		GetQRCodesData    endpoint.Endpoint
	}

	service interface {
		GetDataByQRCodeID(ctx context.Context, id, userID uuid.UUID) (Qrcode, error)
		GetQRCodesData(ctx context.Context, limit, offset int32) ([]Qrcode, error)
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
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
		GetDataByQRCodeID: MakeGetDataByQRCodeIDEndpoint(s),
		GetQRCodesData:    MakeGetQRCodesDataEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetDataByQRCodeID = mdw(e.GetDataByQRCodeID)
			e.GetQRCodesData = mdw(e.GetQRCodesData)
		}
	}

	return e
}

// MakeGetDataByQRCodeIDEndpoint ...
func MakeGetDataByQRCodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		qrcodeUID, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get qrcode id: %w", err)
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.GetDataByQRCodeID(ctx, qrcodeUID, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetQRCodesDataEndpoint ...
func MakeGetQRCodesDataEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetQRCodesData(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
