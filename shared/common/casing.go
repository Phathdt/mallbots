package common

import (
	"encoding/json"
	"strings"
)

// ConvertToSnakeCase converts a string to snake_case
func ConvertToSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}

	return strings.ToLower(result.String())
}

// ConvertMapKeysToSnakeCase recursively converts all keys in a map to snake_case
func ConvertMapKeysToSnakeCase(data interface{}) interface{} {
	switch x := data.(type) {
	case map[string]interface{}:
		newMap := make(map[string]interface{})
		for k, v := range x {
			newMap[ConvertToSnakeCase(k)] = ConvertMapKeysToSnakeCase(v)
		}
		return newMap
	case []interface{}:
		newSlice := make([]interface{}, len(x))
		for i, v := range x {
			newSlice[i] = ConvertMapKeysToSnakeCase(v)
		}
		return newSlice
	default:
		return x
	}
}

// SnakeCaseRawMessage is a wrapper around json.RawMessage that converts keys to snake_case
type SnakeCaseRawMessage json.RawMessage

func (m *SnakeCaseRawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	// Unmarshal the original data
	var data interface{}
	if err := json.Unmarshal([]byte(*m), &data); err != nil {
		return nil, err
	}

	// Convert keys to snake_case
	converted := ConvertMapKeysToSnakeCase(data)

	// Marshal back to JSON
	return json.Marshal(converted)
}

func (m *SnakeCaseRawMessage) UnmarshalJSON(data []byte) error {
	*m = SnakeCaseRawMessage(data)
	return nil
}
