package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalUnstructuredSuccess(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"age": 25,
		"hobbies": ["gaming", "reading", "swimming"],
		"address": {
			"city": "San Francisco",
			"zipcode": 94107
		},
		"is_student": false
	}`)

	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", result["name"].(string))
	assert.Equal(t, 25.0, result["age"].(float64))
	assert.Equal(t, false, result["is_student"].(bool))
	assert.Equal(t, "San Francisco", result["address"].(map[string]interface{})["city"].(string))
	assert.Equal(t, 94107.0, result["address"].(map[string]interface{})["zipcode"].(float64))

	hobbies := result["hobbies"].([]interface{})
	assert.Equal(t, "gaming", hobbies[0].(string))
	assert.Equal(t, "reading", hobbies[1].(string))
	assert.Equal(t, "swimming", hobbies[2].(string))
}

func TestUnmarshalUnstructuredInvalidJSON(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"age": 25,
		"hobbies": ["gaming", "reading", "swimming"],
		"address": {
			"city": "San Francisco",
			"zipcode": 94107,
		"is_student": false
	}`)

	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	assert.NotNil(t, err)
}

func TestUnmarshalUnstructuredMissingFields(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"age": 25,
		"hobbies": ["gaming", "reading", "swimming"]
	}`)

	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", result["name"].(string))
	assert.Equal(t, 25.0, result["age"].(float64))

	hobbies := result["hobbies"].([]interface{})
	assert.Equal(t, "gaming", hobbies[0].(string))
	assert.Equal(t, "reading", hobbies[1].(string))
	assert.Equal(t, "swimming", hobbies[2].(string))

	_, addressExists := result["address"]
	assert.False(t, addressExists)

	_, isStudentExists := result["is_student"]
	assert.False(t, isStudentExists)
}

func TestUnmarshalUnstructuredEmptyJSON(t *testing.T) {
	jsonData := []byte(`{}`)

	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestUnmarshalUnstructuredNestedObjects(t *testing.T) {
	jsonData := []byte(`{
		"name": "Alice",
		"details": {
			"age": 25,
			"hobbies": ["gaming", "reading", "swimming"],
			"address": {
				"city": "San Francisco",
				"zipcode": 94107
			},
			"is_student": false
		}
	}`)

	var result map[string]interface{}
	err := json.Unmarshal(jsonData, &result)
	assert.Nil(t, err)
	assert.Equal(t, "Alice", result["name"].(string))

	details := result["details"].(map[string]interface{})
	assert.Equal(t, 25.0, details["age"].(float64))
	assert.Equal(t, false, details["is_student"].(bool))
	assert.Equal(t, "San Francisco", details["address"].(map[string]interface{})["city"].(string))
	assert.Equal(t, 94107.0, details["address"].(map[string]interface{})["zipcode"].(float64))

	hobbies := details["hobbies"].([]interface{})
	assert.Equal(t, "gaming", hobbies[0].(string))
	assert.Equal(t, "reading", hobbies[1].(string))
	assert.Equal(t, "swimming", hobbies[2].(string))
}
