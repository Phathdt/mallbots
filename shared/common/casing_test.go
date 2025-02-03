package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple camelCase",
			input:    "camelCase",
			expected: "camel_case",
		},
		{
			name:     "multiple uppercase",
			input:    "thisIsATest",
			expected: "this_is_a_test",
		},
		{
			name:     "already snake case",
			input:    "already_snake",
			expected: "already_snake",
		},
		{
			name:     "with numbers",
			input:    "user123Name",
			expected: "user123_name",
		},
		{
			name:     "single letter words",
			input:    "aBC",
			expected: "a_b_c",
		},
		{
			name:     "consecutive uppercase",
			input:    "JSONLD",
			expected: "j_s_o_n_l_d",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConvertMapKeysToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "simple map",
			input: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
			expected: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
		},
		{
			name: "nested map",
			input: map[string]interface{}{
				"userInfo": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
				},
			},
			expected: map[string]interface{}{
				"user_info": map[string]interface{}{
					"first_name": "John",
					"last_name":  "Doe",
				},
			},
		},
		{
			name: "array in map",
			input: map[string]interface{}{
				"userList": []interface{}{
					map[string]interface{}{
						"firstName": "John",
					},
					map[string]interface{}{
						"firstName": "Jane",
					},
				},
			},
			expected: map[string]interface{}{
				"user_list": []interface{}{
					map[string]interface{}{
						"first_name": "John",
					},
					map[string]interface{}{
						"first_name": "Jane",
					},
				},
			},
		},
		{
			name: "mixed types",
			input: map[string]interface{}{
				"userId":   123,
				"userInfo": "someData",
				"userMetadata": map[string]interface{}{
					"lastLogin": "2023-01-01",
				},
			},
			expected: map[string]interface{}{
				"user_id":   123,
				"user_info": "someData",
				"user_metadata": map[string]interface{}{
					"last_login": "2023-01-01",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertMapKeysToSnakeCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSnakeCaseRawMessage(t *testing.T) {
	tests := []struct {
		name           string
		input          interface{}
		expectedOutput string
	}{
		{
			name: "basic conversion",
			input: map[string]interface{}{
				"firstName": "John",
				"lastName":  "Doe",
			},
			expectedOutput: `{"first_name":"John","last_name":"Doe"}`,
		},
		{
			name: "complex nested structure",
			input: map[string]interface{}{
				"userInfo": map[string]interface{}{
					"firstName": "John",
					"lastName":  "Doe",
					"addressInfo": map[string]interface{}{
						"streetName": "Main St",
						"zipCode":    "12345",
					},
				},
				"paymentDetails": []interface{}{
					map[string]interface{}{
						"cardType":   "Credit",
						"cardNumber": "****1234",
					},
				},
			},
			expectedOutput: `{"payment_details":[{"card_type":"Credit","card_number":"****1234"}],"user_info":{"address_info":{"street_name":"Main St","zip_code":"12345"},"first_name":"John","last_name":"Doe"}}`,
		},
		{
			name: "pmm info example",
			input: map[string]interface{}{
				"pmmInfo": map[string]interface{}{
					"sigExpiry":     "1737026779",
					"selectedPmmId": "0x626974666",
					"amountOut":     "333187",
				},
				"transactionHash": "0x338b4969",
			},
			expectedOutput: `{"pmm_info":{"amount_out":"333187","selected_pmm_id":"0x626974666","sig_expiry":"1737026779"},"transaction_hash":"0x338b4969"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert input to JSON
			inputJSON, err := json.Marshal(tt.input)
			assert.NoError(t, err)

			// Create SnakeCaseRawMessage
			var msg SnakeCaseRawMessage
			err = json.Unmarshal(inputJSON, &msg)
			assert.NoError(t, err)

			// Marshal back to JSON
			output, err := json.Marshal(&msg)
			assert.NoError(t, err)

			// Compare with expected output
			assert.JSONEq(t, tt.expectedOutput, string(output))
		})
	}
}

func TestSnakeCaseRawMessage_NilCase(t *testing.T) {
	var msg *SnakeCaseRawMessage
	output, err := json.Marshal(msg)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(output))
}

func TestSnakeCaseRawMessage_InvalidJSON(t *testing.T) {
	// Test with completely invalid JSON
	invalidJSON := []byte(`{"key": invalid`)

	var msg SnakeCaseRawMessage
	err := json.Unmarshal(invalidJSON, &msg)
	assert.Error(t, err, "Unmarshal should fail with invalid JSON")

	// Test with valid JSON but invalid value type
	invalidTypeJSON := []byte(`{"key": [1,2,3`)
	err = json.Unmarshal(invalidTypeJSON, &msg)
	assert.Error(t, err, "Unmarshal should fail with incomplete JSON array")
}
