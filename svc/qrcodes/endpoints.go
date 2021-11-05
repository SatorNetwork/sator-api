package qrcodes

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of qrcode service
	Endpoints struct {
		GetDataByQRCodeID endpoint.Endpoint
		AddQRCode         endpoint.Endpoint
		DeleteQRCodeByID  endpoint.Endpoint
		UpdateQRCode      endpoint.Endpoint
	}

	service interface {
		GetDataByQRCodeID(ctx context.Context, id, userID uuid.UUID) (interface{}, error)
		AddQRCode(ctx context.Context, qr Qrcode) (Qrcode, error)
		DeleteQRCodeByID(ctx context.Context, id uuid.UUID) error
		UpdateQRCode(ctx context.Context, qr Qrcode) error
	}

	// AddQRCodeRequest struct
	AddQRCodeRequest struct {
		ShowID       string  `json:"show_id" validate:"required,uuid"`
		EpisodeID    string  `json:"episode_id" validate:"required,uuid"`
		StartsAt     string  `json:"starts_at" validate:"required"`
		ExpiresAt    string  `json:"expires_at" validate:"required"`
		RewardAmount float64 `json:"reward_amount" validate:"required"`
	}

	// UpdateQRCodeRequest struct
	UpdateQRCodeRequest struct {
		ID           string  `json:"id" validate:"required,uuid"`
		ShowID       string  `json:"show_id" validate:"required,uuid"`
		EpisodeID    string  `json:"episode_id" validate:"required,uuid"`
		StartsAt     string  `json:"starts_at" validate:"required"`
		ExpiresAt    string  `json:"expires_at" validate:"required"`
		RewardAmount float64 `json:"reward_amount" validate:"required"`
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetDataByQRCodeID: MakeGetDataByQRCodeIDEndpoint(s),
		AddQRCode:         MakeAddQRCodeEndpoint(s, validateFunc),
		DeleteQRCodeByID:  MakeDeleteQRCodeByIDEndpoint(s),
		UpdateQRCode:      MakeUpdateQRCodeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetDataByQRCodeID = mdw(e.GetDataByQRCodeID)
			e.AddQRCode = mdw(e.AddQRCode)
			e.DeleteQRCodeByID = mdw(e.DeleteQRCodeByID)
			e.UpdateQRCode = mdw(e.UpdateQRCode)
		}
	}

	return e
}

// MakeGetDataByQRCodeIDEndpoint ...
func MakeGetDataByQRCodeIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		qrcodeUID, err := uuid.Parse(req.(string))
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

// MakeAddQRCodeEndpoint ...
func MakeAddQRCodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(AddQRCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		startsAt, err := utils.DateFromString(req.StartsAt)
		if err != nil {
			return nil, fmt.Errorf("could not add parse start date from string: %w", err)
		}

		expiresAt, err := utils.DateFromString(req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("could not add parse expire date from string: %w", err)
		}

		resp, err := s.AddQRCode(ctx, Qrcode{
			ShowID:       showID,
			EpisodeID:    episodeID,
			StartsAt:     startsAt,
			ExpiresAt:    expiresAt,
			RewardAmount: req.RewardAmount,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeDeleteQRCodeByIDEndpoint ...
func MakeDeleteQRCodeByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get qrcode id: %w", err)
		}

		err = s.DeleteQRCodeByID(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("%w qrcode id: %v", ErrInvalidParameter, err)
		}

		return true, nil
	}
}

// MakeUpdateQRCodeEndpoint ...
func MakeUpdateQRCodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(UpdateQRCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get qrcode id: %w", err)
		}

		showID, err := uuid.Parse(req.ShowID)
		if err != nil {
			return nil, fmt.Errorf("could not get show id: %w", err)
		}

		episodeID, err := uuid.Parse(req.EpisodeID)
		if err != nil {
			return nil, fmt.Errorf("could not get episode id: %w", err)
		}

		startsAt, err := utils.DateFromString(req.StartsAt)
		if err != nil {
			return nil, fmt.Errorf("could not add parse start date from string: %w", err)
		}

		expiresAt, err := utils.DateFromString(req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("could not add parse expire date from string: %w", err)
		}

		err = s.UpdateQRCode(ctx, Qrcode{
			ID:           id,
			ShowID:       showID,
			EpisodeID:    episodeID,
			StartsAt:     startsAt,
			ExpiresAt:    expiresAt,
			RewardAmount: req.RewardAmount,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}
