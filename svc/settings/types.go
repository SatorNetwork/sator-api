package settings

import (
	"fmt"
	"time"

	"github.com/SatorNetwork/sator-api/svc/settings/repository"
)

const (
	SettingPuzzleGamePaidStepsKey = "puzzle_game_paid_steps"
	SettingPuzzleGameRewardsKey   = "puzzle_game_rewards"
)

// castToSetting casts the database record to the setting structure
func castToSetting(rawSetting *repository.Setting) *Setting {
	return &Setting{
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
func stringToSettingsValue(value string, valueType repository.SettingsValueType) interface{} {
	switch valueType {
	case repository.SettingsValueTypeBool:
		return stringToBool(value)
	case repository.SettingsValueTypeFloat:
		return stringToFloat64(value)
	case repository.SettingsValueTypeInt:
		return stringToInt(value)
	case repository.SettingsValueTypeString,
		repository.SettingsValueTypeJson,
		repository.SettingsValueTypeDuration,
		repository.SettingsValueTypeDatetime:
		return value
	default:
		return value
	}
}
