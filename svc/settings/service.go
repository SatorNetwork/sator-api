package settings

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/SatorNetwork/sator-api/svc/settings/repository"
)

type (
	Service struct {
		repo settingsRepository
	}

	settingsRepository interface {
		AddSetting(ctx context.Context, arg repository.AddSettingParams) (repository.Setting, error)
		DeleteSetting(ctx context.Context, key string) error
		GetSettingByKey(ctx context.Context, key string) (repository.Setting, error)
		GetSettings(ctx context.Context) ([]repository.Setting, error)
		UpdateSetting(ctx context.Context, arg repository.UpdateSettingParams) (repository.Setting, error)
	}

	// Setting is the settings for the unity game API
	Setting struct {
		Key         string      `json:"key"`
		Name        string      `json:"name"`
		Value       interface{} `json:"value"`
		ValueType   string      `json:"value_type"`
		Description string      `json:"description,omitempty"`
	}

	// SettingValue represents the setting value db model
	SettingValue struct {
		Value interface{} `json:"value"`
	}
)

func NewService(
	repo settingsRepository,
) *Service {
	return &Service{
		repo: repo,
	}
}

// SettingsValueTypes returns the supported settings value types
func (s *Service) SettingsValueTypes() map[string]string {
	return map[string]string{
		string(repository.SettingsValueTypeBool):   "Boolean",
		string(repository.SettingsValueTypeFloat):  "Float",
		string(repository.SettingsValueTypeInt):    "Integer",
		string(repository.SettingsValueTypeJson):   "JSON",
		string(repository.SettingsValueTypeString): "String",
		// string(repository.SettingsValueTypeDatetime): "DateTime (in RFC3339 format)",
		string(repository.SettingsValueTypeDuration): "Duration (duration string, eg: 5s, 1m, 1h, 2h30m)",
	}
}

// AddSetting adds the setting
func (s *Service) AddSetting(ctx context.Context, key, name, valueType string, value interface{}, description string) (*Setting, error) {
	if name == "" {
		name = key
	}

	res, err := s.repo.AddSetting(ctx, repository.AddSettingParams{
		Key:         toSnakeCase(key),
		Name:        toTitle(name),
		Value:       settingsValueToString(value),
		ValueType:   repository.SettingsValueType(valueType),
		Description: sql.NullString{String: description, Valid: len(description) > 0},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to add setting: %w", err)
	}

	return castToSetting(&res), nil
}

func (s *Service) GetSettingByKey(ctx context.Context, key string) (*Setting, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}
		return nil, err
	}

	return castToSetting(&setting), nil
}

// GetSettings returns the all settings
func (s *Service) GetSettings(ctx context.Context) []*Setting {
	settings, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil
	}

	var result []*Setting
	for _, setting := range settings {
		result = append(result, castToSetting(&setting))
	}

	return result
}

// UpdateSetting updates the setting by key
func (s *Service) UpdateSetting(ctx context.Context, key string, value interface{}) (*Setting, error) {
	res, err := s.repo.UpdateSetting(ctx, repository.UpdateSettingParams{
		Key:   key,
		Value: settingsValueToString(value),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update setting: %w", err)
	}

	return castToSetting(&res), nil
}

// DeleteSetting deletes the setting by key
func (s *Service) DeleteSetting(ctx context.Context, key string) error {
	return s.repo.DeleteSetting(ctx, key)
}

func (s *Service) getValueByKey(ctx context.Context, key string, valueType repository.SettingsValueType) (string, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != valueType {
		return "", fmt.Errorf("key %s value type is not integer, it's %s", key, setting.ValueType)
	}

	return setting.Value, nil
}

// GetBool returns the setting value
func (s *Service) GetBool(ctx context.Context, key string) (bool, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeBool)
	if err != nil {
		return false, err
	}

	return stringToBool(value), nil
}

// GetString returns the setting value
func (s *Service) GetString(ctx context.Context, key string) (string, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeString)
	if err != nil {
		return "", err
	}

	return value, nil
}

// GetFloat64 returns the setting value
func (s *Service) GetFloat64(ctx context.Context, key string) (float64, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeFloat)
	if err != nil {
		return 0, err
	}

	return stringToFloat64(value), nil
}

// GetInt returns the setting value
func (s *Service) GetInt(ctx context.Context, key string) (int, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeInt)
	if err != nil {
		return 0, err
	}

	return stringToInt(value), nil
}

// GetInt32 returns the setting value
func (s *Service) GetInt32(ctx context.Context, key string) (int32, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeInt)
	if err != nil {
		return 0, err
	}

	return stringToInt32(value), nil
}

// GetJSON returns the setting value
func (s *Service) GetJSON(ctx context.Context, key string, result interface{}) error {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeJson)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return fmt.Errorf("failed to unmarshal setting value: %w", err)
	}

	return nil
}

// GetDuration returns the setting value
func (s *Service) GetDuration(ctx context.Context, key string) (time.Duration, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeDuration)
	if err != nil {
		return 0, err
	}

	return stringToDuration(value), nil
}

// GetDatetime returns the setting value
func (s *Service) GetDatetime(ctx context.Context, key string) (time.Time, error) {
	value, err := s.getValueByKey(ctx, key, repository.SettingsValueTypeDatetime)
	if err != nil {
		return time.Time{}, err
	}

	return stringToTime(value)
}
