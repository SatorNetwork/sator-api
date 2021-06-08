package qrcodes

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of qrcode service
	Endpoints struct {
		GetDataByQRCodeID endpoint.Endpoint
	}

	service interface {
		GetDataByQRCodeID(ctx context.Context, id uuid.UUID) (interface{}, error)
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	e := Endpoints{
		GetDataByQRCodeID: MakeGetDataByQRCodeIDEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetDataByQRCodeID = mdw(e.GetDataByQRCodeID)
		}
	}

	return e
}

// MakeGetDataByQRCodeIDEndpoint ...
func MakeGetDataByQRCodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.QRCodeUUIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get qrcode id: %w", err)
		}

		resp, err := s.GetDataByQRCodeID(ctx, uid)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
