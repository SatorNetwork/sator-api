package settings

import (
	"context"
	"github.com/SatorNetwork/sator-api/lib/rbac"
	"github.com/SatorNetwork/sator-api/lib/validator"
	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoints collection of settings service
	Endpoints struct {
		GetSettingsValueTypes endpoint.Endpoint
		GetSettings           endpoint.Endpoint
		GetSettingsByKey      endpoint.Endpoint
		AddSetting            endpoint.Endpoint
		UpdateSetting         endpoint.Endpoint
		DeleteSetting         endpoint.Endpoint
	}

	service interface {
		AddSetting(ctx context.Context, key, name, valueType string, value interface{}, description string) (*Setting, error)
		GetSettingByKey(ctx context.Context, key string) (*Setting, error)
		GetSettings(ctx context.Context) []*Setting
		UpdateSetting(ctx context.Context, key string, value interface{}) (*Setting, error)
		DeleteSetting(ctx context.Context, key string) error
		SettingsValueTypes() map[string]string
	}
)

func MakeEndpoints(s service, m ...endpoint.Middleware) Endpoints {
	validateFunc := validator.ValidateStruct()

	e := Endpoints{
		GetSettings:           MakeGetSettingsEndpoint(s),
		GetSettingsByKey:      MakeGetSettingsByKeyEndpoint(s),
		AddSetting:            MakeAddSettingEndpoint(s, validateFunc),
		UpdateSetting:         MakeUpdateSettingEndpoint(s, validateFunc),
		DeleteSetting:         MakeDeleteSettingEndpoint(s),
		GetSettingsValueTypes: MakeGetSettingsValueTypesEndpoint(s),
	}

	// setup middlewares for each endpoints
	if len(m) > 0 {
		for _, mdw := range m {
			e.GetSettings = mdw(e.GetSettings)
			e.GetSettingsByKey = mdw(e.GetSettingsByKey)
			e.AddSetting = mdw(e.AddSetting)
			e.UpdateSetting = mdw(e.UpdateSetting)
			e.DeleteSetting = mdw(e.DeleteSetting)
			e.GetSettingsValueTypes = mdw(e.GetSettingsValueTypes)
		}
	}

	return e
}

// MakeGetSettingsEndpoint ...
func MakeGetSettingsEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		return s.GetSettings(ctx), nil
	}
}

// AddGameSettingsRequest ...
type AddGameSettingsRequest struct {
	Key         string      `json:"key" validate:"required"`
	Name        string      `json:"name"`
	ValueType   string      `json:"value_type" validate:"required,oneof=int float string json bool duration datetime"`
	Value       interface{} `json:"value" validate:"required"`
	Description string      `json:"description,omitempty"`
}

// MakeAddSettingEndpoint ...
func MakeAddSettingEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		req := request.(AddGameSettingsRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.AddSetting(ctx, req.Key, req.Name, req.ValueType, req.Value, req.Description)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// UpdateGameSettingRequest ...
type UpdateGameSettingRequest struct {
	Key   string      `json:"key" validate:"required"`
	Value interface{} `json:"value" validate:"required"`
}

// MakeUpdateSettingEndpoint ...
func MakeUpdateSettingEndpoint(s service, validateFunc validator.ValidateFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		req := request.(UpdateGameSettingRequest)
		if err := validateFunc(req); err != nil {
			return nil, err
		}

		res, err := s.UpdateSetting(ctx, req.Key, req.Value)
		if err != nil {
			return nil, err
		}

		return res, nil
	}
}

// MakeDeleteSettingEndpoint ...
func MakeDeleteSettingEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		if err := s.DeleteSetting(ctx, request.(string)); err != nil {
			return nil, err
		}

		return true, nil
	}
}

// MakeGetSettingsValueTypesEndpoint ...
func MakeGetSettingsValueTypesEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.RoleAdmin); err != nil {
			return nil, err
		}

		return s.SettingsValueTypes(), nil
	}
}

// MakeGetSettingsByKeyEndpoint ...
func MakeGetSettingsByKeyEndpoint(s service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		if err := rbac.CheckRoleFromContext(ctx, rbac.AvailableForAuthorizedUsers); err != nil {
			return nil, err
		}

		return s.GetSettingByKey(ctx, request.(string))
	}
}
