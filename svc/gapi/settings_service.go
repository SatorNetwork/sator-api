package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/SatorNetwork/sator-api/svc/gapi/repository"
)

type (
	// SettingsService handles communication with the settings related
	SettingsService struct {
		repo settingsRepository
	}

	// settingsRepository interface abstracts the settings storage
	settingsRepository interface {
		AddSetting(ctx context.Context, arg repository.AddSettingParams) (repository.UnityGameSetting, error)
		DeleteSetting(ctx context.Context, key string) error
		GetSettingByKey(ctx context.Context, key string) (repository.UnityGameSetting, error)
		GetSettings(ctx context.Context) ([]repository.UnityGameSetting, error)
		UpdateSetting(ctx context.Context, arg repository.UpdateSettingParams) (repository.UnityGameSetting, error)
	}

	// Settings is the settings for the unity game API
	Settings struct {
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

// NewSettingsService creates a new settings service
func NewSettingsService(repo settingsRepository) *SettingsService {
	return &SettingsService{repo}
}

// SettingsValueTypes returns the supported settings value types
func (s *SettingsService) SettingsValueTypes() map[string]string {
	return map[string]string{
		string(repository.UnityGameSettingsValueTypeBool):   "Boolean",
		string(repository.UnityGameSettingsValueTypeFloat):  "Float",
		string(repository.UnityGameSettingsValueTypeInt):    "Integer",
		string(repository.UnityGameSettingsValueTypeJson):   "JSON",
		string(repository.UnityGameSettingsValueTypeString): "String",
	}
}

// Add adds the setting
func (s *SettingsService) Add(ctx context.Context, key, name, valueType string, value interface{}, description string) (Settings, error) {
	jsonValue, err := json.Marshal(SettingValue{
		Value: value,
	})
	if err != nil {
		return Settings{}, fmt.Errorf("failed to marshal setting value: %w", err)
	}

	key = strings.ToLower(key)
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, ".", "_")
	key = strings.ReplaceAll(key, "-", "_")

	res, err := s.repo.AddSetting(ctx, repository.AddSettingParams{
		Key:         key,
		Name:        name,
		Value:       jsonValue,
		ValueType:   repository.UnityGameSettingsValueType(valueType),
		Description: sql.NullString{String: description, Valid: len(description) > 0},
	})
	if err != nil {
		return Settings{}, fmt.Errorf("failed to add setting: %w", err)
	}

	return castUnityGameSettingToSetting(res), nil
}

// Get function returns the setting by key
func (s *SettingsService) Get(ctx context.Context, key string) (Settings, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return Settings{}, err
	}

	return castUnityGameSettingToSetting(setting), nil
}

// GetValue returns the setting value
func (s *SettingsService) GetValue(ctx context.Context, key string) (interface{}, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	var value SettingValue
	if err := json.Unmarshal(setting.Value, &value); err != nil {
		return nil, fmt.Errorf("failed to unmarshal setting value: %w", err)
	}

	v, err := castSettingValueToValue(value.Value, setting.ValueType)
	if err != nil {
		return nil, fmt.Errorf("failed to cast setting value: %w", err)
	}

	return v, nil
}

// GetAll returns the all settings
func (s *SettingsService) GetAll(ctx context.Context) []Settings {
	settings, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil
	}

	var result []Settings
	for _, setting := range settings {
		result = append(result, castUnityGameSettingToSetting(setting))
	}

	return result
}

// UpdateSetting updates the setting by key
func (s *SettingsService) Update(ctx context.Context, key string, value interface{}) (Settings, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return Settings{}, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if valueType := settingsValueType(value); setting.ValueType != valueType {
		return Settings{}, fmt.Errorf("value type %s is not supported for setting %s", valueType, key)
	}

	jsonValue, err := json.Marshal(SettingValue{
		Value: value,
	})
	if err != nil {
		return Settings{}, fmt.Errorf("failed to marshal setting value: %w", err)
	}

	res, err := s.repo.UpdateSetting(ctx, repository.UpdateSettingParams{
		Key:   key,
		Value: jsonValue,
	})
	if err != nil {
		return Settings{}, fmt.Errorf("failed to update setting: %w", err)
	}

	return castUnityGameSettingToSetting(res), nil
}

// Delete deletes the setting by key
func (s *SettingsService) Delete(ctx context.Context, key string) error {
	return s.repo.DeleteSetting(ctx, key)
}

// settingsValueType returns the value type of the setting
func settingsValueType(v interface{}) repository.UnityGameSettingsValueType {
	switch v.(type) {
	case bool:
		return repository.UnityGameSettingsValueTypeBool
	case float64:
		return repository.UnityGameSettingsValueTypeFloat
	case int:
		return repository.UnityGameSettingsValueTypeInt
	case string:
		return repository.UnityGameSettingsValueTypeString
	case map[string]interface{}:
		return repository.UnityGameSettingsValueTypeJson
	default:
		return repository.UnityGameSettingsValueTypeString
	}
}

// castUnityGameSettingToSetting casts the database record to the setting structure
func castUnityGameSettingToSetting(rawSetting repository.UnityGameSetting) Settings {
	var valRes interface{}
	v := SettingValue{}
	if err := json.Unmarshal(rawSetting.Value, &v); err == nil {
		valRes = v.Value
	} else {
		valRes = string(rawSetting.Value)
	}

	return Settings{
		Key:         rawSetting.Key,
		Name:        rawSetting.Name,
		Value:       valRes,
		ValueType:   string(rawSetting.ValueType),
		Description: rawSetting.Description.String,
	}
}

// cast settings value to the value_type data type
func castSettingValueToValue(value interface{}, valueType repository.UnityGameSettingsValueType) (interface{}, error) {
	switch valueType {
	case repository.UnityGameSettingsValueTypeBool:
		return value.(bool), nil
	case repository.UnityGameSettingsValueTypeFloat:
		return value.(float64), nil
	case repository.UnityGameSettingsValueTypeInt:
		return value.(int), nil
	case repository.UnityGameSettingsValueTypeJson:
		return value.(map[string]interface{}), nil
	case repository.UnityGameSettingsValueTypeString:
		return value.(string), nil
	default:
		return nil, fmt.Errorf("value type %s is not supported", valueType)
	}
}
