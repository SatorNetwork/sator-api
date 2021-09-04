package referrals

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/validator"

	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type (
	// Endpoints collection of profile service
	Endpoints struct {
		AddReferralCodeData        endpoint.Endpoint
		DeleteReferralCodeDataByID endpoint.Endpoint
		GetMyReferralCode          endpoint.Endpoint
		GetReferralCodesDataList   endpoint.Endpoint
		UpdateReferralCodeData     endpoint.Endpoint

		GetReferralsWithPaginationByUserID endpoint.Endpoint
		StoreUserWithValidCode             endpoint.Endpoint
	}

	service interface {
		// Referral codes
		AddReferralCodeData(ctx context.Context, rc ReferralCode) (ReferralCode, error)
		DeleteReferralCodeDataByID(ctx context.Context, id uuid.UUID) error
		GetMyReferralCode(ctx context.Context, uid uuid.UUID, username string) ([]ReferralCode, error)
		GetReferralCodesDataList(ctx context.Context, limit, offset int32) ([]ReferralCode, error)
		UpdateReferralCodeData(ctx context.Context, rc ReferralCode) error

		// Referrals
		GetReferralsWithPaginationByUserID(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]Referral, error)
		StoreUserWithValidCode(ctx context.Context, uid uuid.UUID, code string) (bool, error)
	}

	// PaginationRequest struct
	PaginationRequest struct {
		Page         int32 `json:"page,omitempty" validate:"number,gte=0"`
		ItemsPerPage int32 `json:"items_per_page,omitempty" validate:"number,gte=0"`
	}

	// AddReferralCodeRequest struct
	AddReferralCodeRequest struct {
		Title        string `json:"title,omitempty" validate:"required,gt=0"`
		Code         string `json:"code,omitempty" validate:"required,gt=0"`
		ReferralLink string `json:"referral_link"`
		IsPersonal   bool   `json:"is_personal,omitempty"`
	}

	// UpdateReferralCodeRequest struct
	UpdateReferralCodeRequest struct {
		ID           string `json:"id,omitempty" validate:"required,uuid"`
		Title        string `json:"title,omitempty" validate:"required"`
		Code         string `json:"code,omitempty" validate:"required"`
		ReferralLink string `json:"referral_link"`
		IsPersonal   bool   `json:"is_personal,omitempty" validate:"required"`
		UserID       string `json:"user_id"`
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
		AddReferralCodeData:        MakeAddReferralCodeDataEndpoint(s, validateFunc),
		DeleteReferralCodeDataByID: MakeDeleteReferralCodeDataByIDEndpoint(s),
		GetMyReferralCode:          MakeGetMyReferralCodeEndpoint(s),
		GetReferralCodesDataList:   MakeGetReferralCodesDataListEndpoint(s, validateFunc),
		UpdateReferralCodeData:     MakeUpdateReferralCodeDataEndpoint(s, validateFunc),

		GetReferralsWithPaginationByUserID: MakeGetReferralsWithPaginationByUserIDEndpoint(s, validateFunc),
		StoreUserWithValidCode:             MakeStoreUserWithValidCodeEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.AddReferralCodeData = mdw(e.AddReferralCodeData)
			e.DeleteReferralCodeDataByID = mdw(e.DeleteReferralCodeDataByID)
			e.GetMyReferralCode = mdw(e.GetMyReferralCode)
			e.GetReferralCodesDataList = mdw(e.GetReferralCodesDataList)
			e.UpdateReferralCodeData = mdw(e.UpdateReferralCodeData)

			e.GetReferralsWithPaginationByUserID = mdw(e.GetReferralsWithPaginationByUserID)
			e.StoreUserWithValidCode = mdw(e.StoreUserWithValidCode)
		}
	}

	return e
}

// MakeAddReferralCodeDataEndpoint ...
func MakeAddReferralCodeDataEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(AddReferralCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.AddReferralCodeData(ctx, ReferralCode{
			Title:      req.Title,
			Code:       req.Code,
			IsPersonal: req.IsPersonal,
			UserID:     uid,
		})
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeUpdateReferralCodeDataEndpoint ...
func MakeUpdateReferralCodeDataEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(UpdateReferralCodeRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		id, err := uuid.Parse(req.ID)
		if err != nil {
			return nil, fmt.Errorf("could not get referral code id: %w", err)
		}

		err = s.UpdateReferralCodeData(ctx, ReferralCode{
			ID:           id,
			Title:        req.Title,
			Code:         req.Code,
			ReferralLink: req.ReferralLink,
			IsPersonal:   req.IsPersonal,
			UserID:       uid,
		})
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeDeleteReferralCodeDataByIDEndpoint ...
func MakeDeleteReferralCodeDataByIDEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, err := uuid.Parse(request.(string))
		if err != nil {
			return nil, fmt.Errorf("could not get referral code id: %w", err)
		}

		err = s.DeleteReferralCodeDataByID(ctx, id)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetReferralCodesDataListEndpoint ...
func MakeGetReferralCodesDataListEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetReferralCodesDataList(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetMyReferralCodeEndpoint ...
func MakeGetMyReferralCodeEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		username, err := jwt.UsernameFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get username: %w", err)
		}

		resp, err := s.GetMyReferralCode(ctx, uid, username)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeStoreUserWithValidCodeEndpoint ...
func MakeStoreUserWithValidCodeEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		resp, err := s.StoreUserWithValidCode(ctx, uid, request.(string))
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

// MakeGetReferralsWithPaginationByUserIDEndpoint ...
func MakeGetReferralsWithPaginationByUserIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, err := s.GetReferralsWithPaginationByUserID(ctx, uid, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
