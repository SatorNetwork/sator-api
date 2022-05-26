package gapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
		string(repository.UnityGameSettingsValueTypeBool):     "Boolean",
		string(repository.UnityGameSettingsValueTypeFloat):    "Float",
		string(repository.UnityGameSettingsValueTypeInt):      "Integer",
		string(repository.UnityGameSettingsValueTypeJson):     "JSON",
		string(repository.UnityGameSettingsValueTypeString):   "String",
		string(repository.UnityGameSettingsValueTypeDatetime): "DateTime (in RFC3339 format)",
		string(repository.UnityGameSettingsValueTypeDuration): "Duration (duration string, eg: 5s, 1m, 1h, 2h30m)",
	}
}

// Add adds the setting
func (s *SettingsService) Add(ctx context.Context, key, name, valueType string, value interface{}, description string) (Settings, error) {
	if name == "" {
		name = key
	}
	name = strings.TrimSpace(name)
	name = strings.ToTitle(name)

	res, err := s.repo.AddSetting(ctx, repository.AddSettingParams{
		Key:         toSnakeCase(key),
		Name:        name,
		Value:       settingsValueToString(value),
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
	res, err := s.repo.UpdateSetting(ctx, repository.UpdateSettingParams{
		Key:   key,
		Value: settingsValueToString(value),
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

// GetBool returns the setting value
func (s *SettingsService) GetBool(ctx context.Context, key string) (bool, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return false, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeBool {
		return false, fmt.Errorf("key %s value type is not boolean, it's %s", key, setting.ValueType)
	}

	return stringToBool(setting.Value), nil
}

// GetString returns the setting value
func (s *SettingsService) GetString(ctx context.Context, key string) (string, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeString {
		return "", fmt.Errorf("key %s value type is not string, it's %s", key, setting.ValueType)
	}

	return setting.Value, nil
}

// GetFloat64 returns the setting value
func (s *SettingsService) GetFloat64(ctx context.Context, key string) (float64, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeFloat {
		return 0, fmt.Errorf("key %s value type is not float, it's %s", key, setting.ValueType)
	}

	return stringToFloat64(setting.Value), nil
}

// GetInt returns the setting value
func (s *SettingsService) GetInt(ctx context.Context, key string) (int, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeInt {
		return 0, fmt.Errorf("key %s value type is not integer, it's %s", key, setting.ValueType)
	}

	return stringToInt(setting.Value), nil
}

// GetJSON returns the setting value
func (s *SettingsService) GetJSON(ctx context.Context, key string, result interface{}) error {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeJson {
		return fmt.Errorf("key %s value type is not json, it's %s", key, setting.ValueType)
	}

	if err := json.Unmarshal([]byte(setting.Value), &result); err != nil {
		return fmt.Errorf("failed to unmarshal setting value: %w", err)
	}

	return nil
}

// GetDurration returns the setting value
func (s *SettingsService) GetDurration(ctx context.Context, key string) (time.Duration, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeString {
		return 0, fmt.Errorf("key %s value type is not string, it's %s", key, setting.ValueType)
	}

	return stringToDuration(setting.Value), nil
}

// GetDatetime returns the setting value
func (s *SettingsService) GetDatetime(ctx context.Context, key string) (time.Time, error) {
	setting, err := s.repo.GetSettingByKey(ctx, key)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get setting with key %s: %w", key, err)
	}

	if setting.ValueType != repository.UnityGameSettingsValueTypeString {
		return time.Time{}, fmt.Errorf("key %s value type is not string, it's %s", key, setting.ValueType)
	}

	return stringToTime(setting.Value)
}

// castUnityGameSettingToSetting casts the database record to the setting structure
func castUnityGameSettingToSetting(rawSetting repository.UnityGameSetting) Settings {
	return Settings{
		Key:         rawSetting.Key,
		Name:        rawSetting.Name,
		Value:       stringToSettingsValue(rawSetting.Value, rawSetting.ValueType),
		ValueType:   string(rawSetting.ValueType),
		Description: rawSetting.Description.String,
	}
}

// convert settings value to string
func settingsValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return boolToString(v)
	case float64:
		return float64ToString(v)
	case int:
		return intToString(v)
	case map[string]interface{}:
		return mapToString(v)
	case time.Duration:
		return durationToString(v)
	case time.Time:
		return timeToString(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// convert string to settings value
func stringToSettingsValue(value string, valueType repository.UnityGameSettingsValueType) interface{} {
	switch valueType {
	case repository.UnityGameSettingsValueTypeBool:
		return stringToBool(value)
	case repository.UnityGameSettingsValueTypeFloat:
		return stringToFloat64(value)
	case repository.UnityGameSettingsValueTypeInt:
		return stringToInt(value)
	case repository.UnityGameSettingsValueTypeString:
		return value
	case repository.UnityGameSettingsValueTypeJson:
		return stringToMap(value)
	case repository.UnityGameSettingsValueTypeDuration:
		return stringToDuration(value)
	case repository.UnityGameSettingsValueTypeDatetime:
		res, _ := stringToTime(value)
		return res
	default:
		return value
	}
}
