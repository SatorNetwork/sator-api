package settings

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// convert float64 to string
func float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// convert string to float64
func stringToFloat64(s string) float64 {
	res, _ := strconv.ParseFloat(s, 64)
	return res
}

// convert int to string
func intToString(i int) string {
	return strconv.Itoa(i)
}

// convert string to int
func stringToInt(s string) int {
	res, _ := strconv.Atoi(s)
	return res
}

// convert string to int
func stringToInt32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

// convert bool to string
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// convert string to bool
func stringToBool(s string) bool {
	switch s {
	case "true", "1", "yes", "on", "enable", "enabled", "active", "ok", "success", "pass", "passed", "t", "y":
		return true
	case "false", "0", "no", "off", "disable", "disabled", "inactive", "fail", "failed", "f", "n", "nopass", "nopassed":
		return false
	default:
		return false
	}
}

// convert time.Duration to string
func durationToString(d time.Duration) string {
	return d.String()
}

// convert string to time.Duration
func stringToDuration(s string) time.Duration {
	d, _ := time.ParseDuration(s)
	return d
}

// convert map to json string
func mapToString(m map[string]interface{}) string {
	json, _ := json.Marshal(m)
	return string(json)
}

// convert json string to map
// func stringToMap(s string) map[string]interface{} {
// 	var m map[string]interface{}
// 	json.Unmarshal([]byte(s), &m)
// 	return m
// }

// convert time.Time to string
func timeToString(t time.Time) string {
	return t.Format(time.RFC3339)
}

// convert string to time.Time
func stringToTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// format string to snake case
func toSnakeCase(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = alphaNumUnderscore(s)

	return s
}

// get only alphanumeric characters from string and replace spaces with underscores
func alphaNumUnderscore(s string) string {
	return strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		if r == ' ' || r == '_' || r == '-' || r == '.' || r == ',' || r == ':' || r == ';' || r == '(' || r == ')' || r == '[' || r == ']' || r == '{' || r == '}' {
			return '_'
		}
		return -1
	}, s)
}

// format string to title case
// e.g. "my_string" -> "My String"
func toTitle(s string) string {
	return cases.
		Title(language.Und).
		String(strings.TrimSpace(strings.ReplaceAll(s, "_", " ")))
}
