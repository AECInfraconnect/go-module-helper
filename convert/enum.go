// Package convert provides utilities for enum value conversion and validation.
// This package helps with converting between enum types, validating enum values,
// and performing flexible enum matching operations.
package convert

import (
	"fmt"
	"strings"
)

// EnumToString converts an enum value to its string representation.
//
// Example:
//
//	type Status int
//	const (
//	    Active Status = 1
//	    Inactive Status = 0
//	)
//	str := convert.EnumToString(Active) // Returns "1"
func EnumToString(enum interface{}) string {
	return ToString(enum)
}

// StringToEnum converts a string to an enum value using a provided mapping.
// Returns an error if the string value is not found in the enum map.
//
// Example:
//
//	statusMap := map[string]interface{}{
//	    "ACTIVE": 1,
//	    "INACTIVE": 0,
//	}
//	value, err := convert.StringToEnum("ACTIVE", statusMap)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	status := value.(int) // Returns 1
func StringToEnum(str string, enumMap map[string]interface{}) (interface{}, error) {
	if value, ok := enumMap[str]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("invalid enum value: %s", str)
}

// EnumToInt converts an enum value to an integer.
// Returns an error if the conversion is not possible.
//
// Example:
//
//	type Priority string
//	const High Priority = "HIGH"
//	num, err := convert.EnumToInt(High) // Error: cannot convert string to int
func EnumToInt(enum interface{}) (int, error) {
	return ToInt(enum)
}

// IntToEnum converts an integer to an enum value using a provided mapping.
// Returns an error if the integer value is not found in the enum map.
//
// Example:
//
//	statusMap := map[int]interface{}{
//	    1: "ACTIVE",
//	    0: "INACTIVE",
//	}
//	value, err := convert.IntToEnum(1, statusMap)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	status := value.(string) // Returns "ACTIVE"
func IntToEnum(num int, enumMap map[int]interface{}) (interface{}, error) {
	if value, ok := enumMap[num]; ok {
		return value, nil
	}
	return nil, fmt.Errorf("invalid enum value: %d", num)
}

// EnumExists checks if a string value exists as a key in the enum map.
//
// Example:
//
//	statusMap := map[string]interface{}{
//	    "ACTIVE": 1,
//	    "INACTIVE": 0,
//	}
//	exists := convert.EnumExists("ACTIVE", statusMap) // Returns true
//	exists = convert.EnumExists("PENDING", statusMap) // Returns false
func EnumExists(value string, enumMap map[string]interface{}) bool {
	_, ok := enumMap[value]
	return ok
}

// EnumExistsInt checks if an integer value exists as a key in the enum map.
//
// Example:
//
//	statusMap := map[int]interface{}{
//	    1: "ACTIVE",
//	    0: "INACTIVE",
//	}
//	exists := convert.EnumExistsInt(1, statusMap) // Returns true
//	exists = convert.EnumExistsInt(2, statusMap) // Returns false
func EnumExistsInt(value int, enumMap map[int]interface{}) bool {
	_, ok := enumMap[value]
	return ok
}

// GetEnumKeys extracts and returns all keys from an enum map as a string slice.
// The order of keys is not guaranteed.
//
// Example:
//
//	statusMap := map[string]interface{}{
//	    "ACTIVE": 1,
//	    "INACTIVE": 0,
//	}
//	keys := convert.GetEnumKeys(statusMap) // Returns ["ACTIVE", "INACTIVE"] (order may vary)
func GetEnumKeys(enumMap map[string]interface{}) []string {
	keys := make([]string, 0, len(enumMap))
	for k := range enumMap {
		keys = append(keys, k)
	}
	return keys
}

// GetEnumValues extracts and returns all values from an enum map as an interface slice.
// The order of values is not guaranteed.
//
// Example:
//
//	statusMap := map[string]interface{}{
//	    "ACTIVE": 1,
//	    "INACTIVE": 0,
//	}
//	values := convert.GetEnumValues(statusMap) // Returns [1, 0] (order may vary)
func GetEnumValues(enumMap map[string]interface{}) []interface{} {
	values := make([]interface{}, 0, len(enumMap))
	for _, v := range enumMap {
		values = append(values, v)
	}
	return values
}

// NormalizeEnumString normalizes a string for enum comparison by converting it to
// uppercase and trimming surrounding whitespace. Useful for case-insensitive matching.
//
// Example:
//
//	normalized := convert.NormalizeEnumString("  active  ") // Returns "ACTIVE"
//	normalized = convert.NormalizeEnumString("InActive")   // Returns "INACTIVE"
func NormalizeEnumString(str string) string {
	return strings.ToUpper(strings.TrimSpace(str))
}

// EnumMatch performs case-insensitive enum matching by normalizing both the input
// string and enum map keys before comparison. Returns the matched enum value.
//
// Example:
//
//	statusMap := map[string]interface{}{
//	    "ACTIVE": 1,
//	    "INACTIVE": 0,
//	}
//	value, err := convert.EnumMatch("active", statusMap)   // Returns 1 (case-insensitive match)
//	value, err = convert.EnumMatch("  ACTIVE  ", statusMap) // Returns 1 (whitespace trimmed)
//	value, err = convert.EnumMatch("pending", statusMap)    // Returns error
func EnumMatch(str string, enumMap map[string]interface{}) (interface{}, error) {
	normalized := NormalizeEnumString(str)

	for k, v := range enumMap {
		if NormalizeEnumString(k) == normalized {
			return v, nil
		}
	}

	return nil, fmt.Errorf("invalid enum value: %s", str)
}

// ValidateEnum checks if a value exists in the list of valid enum values.
// Returns the validated value if found, otherwise returns an error with all valid values.
//
// Example:
//
//	validStatuses := []string{"ACTIVE", "INACTIVE", "PENDING"}
//	status, err := convert.ValidateEnum("ACTIVE", validStatuses) // Returns "ACTIVE"
//	status, err = convert.ValidateEnum("DELETED", validStatuses) // Returns error with valid values
func ValidateEnum(value string, validValues []string) (string, error) {
	for _, valid := range validValues {
		if value == valid {
			return value, nil
		}
	}
	return "", fmt.Errorf("invalid enum value: %s, valid values: %v", value, validValues)
}

// ValidateEnumInt checks if an integer value exists in the list of valid enum integers.
// Returns the validated value if found, otherwise returns an error with all valid values.
//
// Example:
//
//	validPriorities := []int{1, 2, 3, 4, 5}
//	priority, err := convert.ValidateEnumInt(3, validPriorities) // Returns 3
//	priority, err = convert.ValidateEnumInt(10, validPriorities) // Returns error with valid values
func ValidateEnumInt(value int, validValues []int) (int, error) {
	for _, valid := range validValues {
		if value == valid {
			return value, nil
		}
	}
	return 0, fmt.Errorf("invalid enum value: %d, valid values: %v", value, validValues)
}
