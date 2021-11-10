package referrals

import (
	"context"
	"fmt"

	"github.com/SatorNetwork/sator-api/internal/httpencoder"
	"github.com/SatorNetwork/sator-api/internal/jwt"
	"github.com/SatorNetwork/sator-api/internal/rbac"
	"github.com/SatorNetwork/sator-api/internal/utils"
	"github.com/SatorNetwork/sator-api/internal/validator"

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
		GetMyReferralCode(ctx context.Context, uid uuid.UUID, username string) (ReferralCode, error)
		GetReferralCodesDataList(ctx context.Context, limit, offset int32) ([]ReferralCode, int64, error)
		UpdateReferralCodeData(ctx context.Context, rc ReferralCode) error

		// Referrals
		GetReferralsWithPaginationByUserID(ctx context.Context, uid uuid.UUID, limit, offset int32) ([]Referral, error)
		StoreUserWithValidCode(ctx context.Context, uid uuid.UUID, code string) (bool, error)
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
		IsPersonal   bool   `json:"is_personal,omitempty"`
		UserID       string `json:"user_id"`
	}
)

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
		// FIXME: is allowed roles correct???
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
			UserID:     &uid,
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
			UserID:       &uid,
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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

		req := request.(utils.PaginationRequest)
		if err := v(req); err != nil {
			return nil, err
		}

		resp, numberOfReferralCodes, err := s.GetReferralCodesDataList(ctx, req.Limit(), req.Offset())
		if err != nil {
			return nil, err
		}

		response := httpencoder.ListResponse{}
		response.Data = resp
		response.Meta.TotalItems = numberOfReferralCodes
		response.Meta.ItemsPerPage = int64(req.ItemsPerPage)
		response.Meta.Page = int64(req.Page)

		return response, nil
	}
}

// MakeGetMyReferralCodeEndpoint ...
func MakeGetMyReferralCodeEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

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
		// FIXME: is allowed roles correct???
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin, rbac.RoleContentManager); err != nil {
			return nil, err
		}

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
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		uid, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user profile id: %w", err)
		}

		req := request.(utils.PaginationRequest)
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
