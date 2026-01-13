// Package convert provides utilities for JSON serialization, deserialization,
// and transformation between different data structures.
// This package simplifies working with JSON data in Go applications.
package convert

import (
	"encoding/json"
	"fmt"
)

// ToJSON converts any Go value to a JSON string representation.
// Returns an error if the value cannot be marshaled to JSON.
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	user := User{Name: "John", Age: 30}
//	jsonStr, err := convert.ToJSON(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Returns: {"name":"John","age":30}
func ToJSON(value interface{}) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// ToJSONIndent converts any Go value to a formatted JSON string with custom indentation.
// The prefix and indent parameters control the output formatting.
//
// Example:
//
//	user := User{Name: "John", Age: 30}
//	jsonStr, err := convert.ToJSONIndent(user, "", "  ")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Returns:
//	// {
//	//   "name": "John",
//	//   "age": 30
//	// }
func ToJSONIndent(value interface{}, prefix, indent string) (string, error) {
	bytes, err := json.MarshalIndent(value, prefix, indent)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FromJSON parses a JSON string and returns the result as interface{}.
// The result can be a map, slice, or primitive value depending on the JSON structure.
//
// Example:
//
//	jsonStr := `{"name":"John","age":30}`
//	result, err := convert.FromJSON(jsonStr)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	data := result.(map[string]interface{})
//	name := data["name"].(string) // "John"
//	age := data["age"].(float64)  // 30
func FromJSON(jsonStr string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// FromJSONTo parses a JSON string into a specific target struct or type.
// The target parameter must be a pointer to the destination variable.
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	jsonStr := `{"name":"John","age":30}`
//	var user User
//	err := convert.FromJSONTo(jsonStr, &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// user.Name = "John", user.Age = 30
func FromJSONTo(jsonStr string, target interface{}) error {
	return json.Unmarshal([]byte(jsonStr), target)
}

// StructToMap converts any struct to a map[string]interface{} representation.
// This is useful for dynamic field access or when working with generic data structures.
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	user := User{Name: "John", Age: 30}
//	dataMap, err := convert.StructToMap(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// dataMap = map[string]interface{}{"name": "John", "age": 30}
//	name := dataMap["name"].(string)
func StructToMap(value interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct: %w", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal to map: %w", err)
	}

	return result, nil
}

// MapToStruct converts a map[string]interface{} to a specific struct type.
// The target parameter must be a pointer to the destination struct.
//
// Example:
//
//	type User struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//	dataMap := map[string]interface{}{
//	    "name": "John",
//	    "age":  30,
//	}
//	var user User
//	err := convert.MapToStruct(dataMap, &user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// user.Name = "John", user.Age = 30
func MapToStruct(data map[string]interface{}, target interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal map: %w", err)
	}

	err = json.Unmarshal(jsonBytes, target)
	if err != nil {
		return fmt.Errorf("failed to unmarshal to struct: %w", err)
	}

	return nil
}

// ToJSONBytes converts any Go value to JSON bytes.
// This is more efficient than ToJSON when you need bytes directly.
//
// Example:
//
//	user := User{Name: "John", Age: 30}
//	jsonBytes, err := convert.ToJSONBytes(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// Use jsonBytes for file writing or HTTP responses
//	ioutil.WriteFile("user.json", jsonBytes, 0644)
func ToJSONBytes(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}

// FromJSONBytes parses JSON bytes and returns the result as interface{}.
// This is more efficient than FromJSON when you already have bytes.
//
// Example:
//
//	jsonBytes := []byte(`{"name":"John","age":30}`)
//	result, err := convert.FromJSONBytes(jsonBytes)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	data := result.(map[string]interface{})
func FromJSONBytes(data []byte) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// IsValidJSON checks whether a string contains valid JSON.
// Returns true if the string can be successfully parsed as JSON, false otherwise.
//
// Example:
//
//	valid := convert.IsValidJSON(`{"name":"John"}`)        // Returns true
//	valid = convert.IsValidJSON(`{name:"John"}`)           // Returns false (invalid syntax)
//	valid = convert.IsValidJSON(`["apple", "banana"]`)     // Returns true
//	valid = convert.IsValidJSON(`not json`)                // Returns false
func IsValidJSON(jsonStr string) bool {
	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}
