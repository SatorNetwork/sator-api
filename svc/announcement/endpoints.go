package announcement

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/SatorNetwork/sator-api/lib/jwt"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"
)

type (
	Endpoints struct {
		CreateAnnouncement      endpoint.Endpoint
		GetAnnouncementByID     endpoint.Endpoint
		UpdateAnnouncement      endpoint.Endpoint
		DeleteAnnouncement      endpoint.Endpoint
		ListAnnouncements       endpoint.Endpoint
		ListUnreadAnnouncements endpoint.Endpoint
		ListActiveAnnouncements endpoint.Endpoint
		MarkAsRead              endpoint.Endpoint
		MarkAllAsRead           endpoint.Endpoint
	}

	service interface {
		CreateAnnouncement(ctx context.Context, req *CreateAnnouncementRequest) (*CreateAnnouncementResponse, error)
		GetAnnouncementByID(ctx context.Context, req *GetAnnouncementByIDRequest) (*Announcement, error)
		UpdateAnnouncementByID(ctx context.Context, req *UpdateAnnouncementRequest) error
		DeleteAnnouncementByID(ctx context.Context, req *DeleteAnnouncementRequest) error
		ListAnnouncements(ctx context.Context) ([]*Announcement, error)
		ListUnreadAnnouncements(ctx context.Context, userID uuid.UUID) ([]*Announcement, error)
		ListActiveAnnouncements(ctx context.Context) ([]*Announcement, error)
		MarkAsRead(ctx context.Context, userID uuid.UUID, req *MarkAsReadRequest) error
		MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		CreateAnnouncement:      MakeCreateAnnouncementEndpoint(s, validateFunc),
		GetAnnouncementByID:     MakeGetAnnouncementByIDEndpoint(s, validateFunc),
		UpdateAnnouncement:      MakeUpdateAnnouncementEndpoint(s, validateFunc),
		DeleteAnnouncement:      MakeDeleteAnnouncementEndpoint(s, validateFunc),
		ListAnnouncements:       MakeListAnnouncementsEndpoint(s, validateFunc),
		ListUnreadAnnouncements: MakeListUnreadAnnouncementsEndpoint(s, validateFunc),
		ListActiveAnnouncements: MakeListActiveAnnouncementsEndpoint(s, validateFunc),
		MarkAsRead:              MakeMarkAsReadEndpoint(s, validateFunc),
		MarkAllAsRead:           MakeMarkAllAsReadEndpoint(s, validateFunc),
	}

	// setup middlewares for each endpoint
	if len(m) > 0 {
		for _, mdw := range m {
			e.CreateAnnouncement = mdw(e.CreateAnnouncement)
			e.GetAnnouncementByID = mdw(e.GetAnnouncementByID)
			e.UpdateAnnouncement = mdw(e.UpdateAnnouncement)
			e.DeleteAnnouncement = mdw(e.DeleteAnnouncement)
			e.ListAnnouncements = mdw(e.ListAnnouncements)
			e.ListUnreadAnnouncements = mdw(e.ListUnreadAnnouncements)
			e.ListActiveAnnouncements = mdw(e.ListActiveAnnouncements)
			e.MarkAsRead = mdw(e.MarkAsRead)
			e.MarkAllAsRead = mdw(e.MarkAllAsRead)
		}
	}

	return e
}

func MakeCreateAnnouncementEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*CreateAnnouncementRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to create-announcement-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		resp, err := s.CreateAnnouncement(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeGetAnnouncementByIDEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*GetAnnouncementByIDRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to get-announcement-by-id")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		resp, err := s.GetAnnouncementByID(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeUpdateAnnouncementEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*UpdateAnnouncementRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to update-announcement-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		err := s.UpdateAnnouncementByID(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeDeleteAnnouncementEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		typedReq, ok := req.(*DeleteAnnouncementRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to delete-announcement-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		err := s.DeleteAnnouncementByID(ctx, typedReq)
		if err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeListAnnouncementsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		resp, err := s.ListAnnouncements(ctx)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeListUnreadAnnouncementsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}
		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		resp, err := s.ListUnreadAnnouncements(ctx, userID)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeListActiveAnnouncementsEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		resp, err := s.ListActiveAnnouncements(ctx)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}

func MakeMarkAsReadEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}
		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		typedReq, ok := req.(*MarkAsReadRequest)
		if !ok {
			return nil, errors.Errorf("can't cast untyped request to mark-as-read-request")
		}
		if err := v(typedReq); err != nil {
			return nil, err
		}

		if err := s.MarkAsRead(ctx, userID, typedReq); err != nil {
			return nil, err
		}

		return true, nil
	}
}

func MakeMarkAllAsReadEndpoint(s service, v validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}
		userID, err := jwt.UserIDFromContext(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get user id: %w", err)
		}

		if err := s.MarkAllAsRead(ctx, userID); err != nil {
			return nil, err
		}

		return true, nil
	}
}
