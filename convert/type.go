// Package convert provides utilities for type conversion between common Go data types.
// This package offers flexible type conversion with automatic type detection and
// support for multiple primitive types including strings, integers, floats, and booleans.
package convert

import (
	"fmt"
	"strconv"
)

// ToString converts any value to its string representation.
// Returns an empty string for nil values.
// Handles string, int, int64, float32, float64, bool, and other types via fmt.Sprintf.
//
// Example:
//
//	str := convert.ToString(123)           // Returns "123"
//	str = convert.ToString(45.67)          // Returns "45.67"
//	str = convert.ToString(true)           // Returns "true"
//	str = convert.ToString("hello")        // Returns "hello"
//	str = convert.ToString(nil)            // Returns ""
func ToString(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts various types to an integer value.
// Supports conversion from int, int64, float32, float64, string, bool, and nil.
// Returns an error if the conversion is not possible.
//
// Example:
//
//	num, err := convert.ToInt("123")       // Returns 123
//	num, err = convert.ToInt(45.67)        // Returns 45 (truncates decimal)
//	num, err = convert.ToInt(true)         // Returns 1
//	num, err = convert.ToInt(false)        // Returns 0
//	num, err = convert.ToInt(nil)          // Returns 0
//	num, err = convert.ToInt("abc")        // Returns error
func ToInt(value interface{}) (int, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int", value)
	}
}

// ToInt64 converts various types to a 64-bit integer value.
// Supports conversion from int64, int, float32, float64, string, bool, and nil.
// Returns an error if the conversion is not possible.
//
// Example:
//
//	num, err := convert.ToInt64("12345678901234")  // Returns 12345678901234
//	num, err = convert.ToInt64(100)                 // Returns 100
//	num, err = convert.ToInt64(99.99)               // Returns 99 (truncates decimal)
//	num, err = convert.ToInt64(true)                // Returns 1
//	num, err = convert.ToInt64(nil)                 // Returns 0
func ToInt64(value interface{}) (int64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case float32:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// ToFloat64 converts various types to a 64-bit floating-point value.
// Supports conversion from float64, float32, int, int64, string, bool, and nil.
// Returns an error if the conversion is not possible.
//
// Example:
//
//	num, err := convert.ToFloat64("123.456")   // Returns 123.456
//	num, err = convert.ToFloat64(100)          // Returns 100.0
//	num, err = convert.ToFloat64(true)         // Returns 1.0
//	num, err = convert.ToFloat64(false)        // Returns 0.0
//	num, err = convert.ToFloat64(nil)          // Returns 0.0
//	num, err = convert.ToFloat64("abc")        // Returns error
func ToFloat64(value interface{}) (float64, error) {
	if value == nil {
		return 0, nil
	}

	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	case bool:
		if v {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// ToBool converts various types to a boolean value.
// Returns false for nil, zero values, empty strings, and "false" strings.
// Returns true for non-zero numbers and "true"/"1"/"t"/"T"/"TRUE" strings.
//
// Example:
//
//	b := convert.ToBool(1)              // Returns true
//	b = convert.ToBool(0)               // Returns false
//	b = convert.ToBool("true")          // Returns true
//	b = convert.ToBool("false")         // Returns false
//	b = convert.ToBool(3.14)            // Returns true
//	b = convert.ToBool(0.0)             // Returns false
//	b = convert.ToBool(nil)             // Returns false
func ToBool(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return v
	case int, int64:
		return v != 0
	case float64, float32:
		return v != 0
	case string:
		b, _ := strconv.ParseBool(v)
		return b
	default:
		return false
	}
}

// ToStringSlice converts any value to a string slice.
// If the input is already a []string, returns it as-is.
// If the input is []interface{}, converts each element to string.
// For any other type, returns a slice containing the stringified value.
//
// Example:
//
//	slice := convert.ToStringSlice([]interface{}{1, 2, 3})      // Returns ["1", "2", "3"]
//	slice = convert.ToStringSlice([]string{"a", "b"})           // Returns ["a", "b"]
//	slice = convert.ToStringSlice("hello")                      // Returns ["hello"]
//	slice = convert.ToStringSlice(nil)                          // Returns []
func ToStringSlice(value interface{}) []string {
	if value == nil {
		return []string{}
	}

	switch v := value.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, len(v))
		for i, item := range v {
			result[i] = ToString(item)
		}
		return result
	default:
		return []string{ToString(value)}
	}
}

// ToIntSlice converts any value to an integer slice.
// If the input is already a []int, returns it as-is.
// If the input is []interface{}, converts each element to int.
// For any other type, attempts to convert it to int and returns a slice with that value.
// Returns an error if any element cannot be converted to int.
//
// Example:
//
//	slice, err := convert.ToIntSlice([]interface{}{"1", "2", "3"})  // Returns [1, 2, 3]
//	slice, err = convert.ToIntSlice([]int{10, 20, 30})              // Returns [10, 20, 30]
//	slice, err = convert.ToIntSlice(42)                             // Returns [42]
//	slice, err = convert.ToIntSlice(nil)                            // Returns []
//	slice, err = convert.ToIntSlice([]interface{}{"a", "b"})        // Returns error
func ToIntSlice(value interface{}) ([]int, error) {
	if value == nil {
		return []int{}, nil
	}

	switch v := value.(type) {
	case []int:
		return v, nil
	case []interface{}:
		result := make([]int, len(v))
		for i, item := range v {
			num, err := ToInt(item)
			if err != nil {
				return nil, err
			}
			result[i] = num
		}
		return result, nil
	default:
		num, err := ToInt(value)
		if err != nil {
			return nil, err
		}
		return []int{num}, nil
	}
}
