package helper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

// ValidateKeyExists checks if all specified keys exist in the params map.
// Returns a map of errors for each missing key, or nil if all keys exist.
func ValidateKeyExists(keys []string, params map[string]interface{}) map[string]error {
	if len(keys) == 0 || params == nil {
		return nil
	}
	var errs = make(map[string]error, 0)

	for _, key := range keys {
		if _, ok := params[key]; !ok {
			message := fmt.Sprintf("key '%s' not exists", key)
			errs[key] = errors.New(message)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

// ValidateNotSpace validates that the value is not a single space character.
// The value must be a string type.
// Returns an error if the value is not a string or is a single space.
func ValidateNotSpace(val interface{}) error {
	if err := ValidateTypeString(val); err != nil {
		return err
	}

	if val.(string) == " " {
		return fmt.Errorf(`can not be space only`)
	}

	return nil
}

// ValidateOnlyThaiLetterNumeric validates that the value contains only Thai letters and digits.
// English letters are not allowed. The value must be a string type.
// Returns an error if the value is not a string or contains invalid characters.
func ValidateOnlyThaiLetterNumeric(val interface{}) error {
	if err := ValidateTypeString(val); err != nil {
		return err
	}

	thaiRegex := `[\x{0E00}-\x{0E7Fa}0-9]`
	engRegex := `[a-zA-Z]`
	exp := regexp.MustCompile(thaiRegex)
	enExp := regexp.MustCompile(engRegex)

	if !exp.MatchString(val.(string)) || enExp.MatchString(val.(string)) {
		return errors.New("letter can be Thai letters and digits only")
	}

	return nil
}

// ValidateUUIDOrIDZero validates that the value is either a valid UUID or the string "0".
// The value must be a string type.
// Returns an error if the value is not a string or is neither a valid UUID nor "0".
func ValidateUUIDOrIDZero(val interface{}) error {
	if err := ValidateTypeString(val); err != nil {
		return err
	}

	if val.(string) == "0" {
		return nil
	}

	if _, err := uuid.Parse(val.(string)); err != nil {
		return errors.New("value is not a valid UUID")
	}

	return nil
}

// ValidateTypeUUID validates that the value is a valid UUID.
// The value must be a string type in UUID format (e.g., "123e4567-e89b-12d3-a456-426614174000").
// Returns an error if the value is not a string or not a valid UUID.
func ValidateTypeUUID(val interface{}) error {
	if err := ValidateTypeString(val); err != nil {
		return err
	}

	if _, err := uuid.Parse(val.(string)); err != nil {
		return errors.New("value is not a valid UUID")
	}

	return nil
}

// ValidateTypeString validates that the value is of type string.
// Returns an error if the value is not a string type.
func ValidateTypeString(val interface{}) error {
	rf := reflect.ValueOf(val)
	if rf.Kind() != reflect.String {
		return fmt.Errorf("value is not type string")
	}
	return nil
}

// ValidateTypeInt validates that the value is of type int.
// Also accepts float64 values that can be safely converted to int without precision loss.
// Returns an error if the value is not an int or convertible float64.
func ValidateTypeInt(val interface{}) error {
	rf := reflect.ValueOf(val)
	if rf.Kind() == reflect.Float64 {
		stringValue := fmt.Sprintf("%d", int(val.(float64)))
		_, err := strconv.ParseInt(stringValue, 10, 32)
		if err == nil {
			return nil
		}
	}
	if rf.Kind() != reflect.Int {
		return fmt.Errorf("value is not type int")
	}
	return nil
}

// ValidateTypeFloat validates that the value is of type float64.
// Returns an error if the value is not a float64 type.
func ValidateTypeFloat(val interface{}) error {
	rf := reflect.ValueOf(val)
	if rf.Kind() != reflect.Float64 {
		return errors.New("value is not type float")
	}
	return nil
}

// ValidateTypeMap validates that the value is of type map.
// Returns an error if the value is not a map type.
func ValidateTypeMap(val interface{}) error {
	rf := reflect.ValueOf(val)
	if rf.Kind() != reflect.Map {
		return errors.New("value is not type map")
	}
	return nil
}

// ValidateTypeSlice validates that the value is of type slice (array).
// Returns an error if the value is not a slice type.
func ValidateTypeSlice(val interface{}) error {
	rf := reflect.ValueOf(val)
	if rf.Kind() != reflect.Slice {
		return errors.New("value is not type array")
	}
	return nil
}

// ValidateTypeBool validates that the value is of type bool.
// Returns an error if the value is not a boolean type.
func ValidateTypeBool(v interface{}) error {
	if reflect.TypeOf(v).Kind() == reflect.Bool {
		return nil
	}
	return errors.New("value is not type bool")
}

// ValidateTypeBoolString validates that the value is either a bool type or a string that can be parsed as bool.
// Accepts: true, false, "true", "false", "1", "0", "t", "f", "T", "F", "TRUE", "FALSE", "True", "False".
// Returns an error if the value is neither a bool nor a parseable bool string.
func ValidateTypeBoolString(v interface{}) error {
	if reflect.TypeOf(v).Kind() == reflect.Bool {
		return nil
	}
	if reflect.TypeOf(v).Kind() == reflect.String {
		if _, err := strconv.ParseBool(v.(string)); err == nil {
			return nil
		}
	}
	return errors.New("value is not type bool")
}

// ValidateTypeMapWithNull validates that the value is either nil or of type map.
// Returns an error if the value is not nil and not a map type.
func ValidateTypeMapWithNull(val interface{}) error {
	if val == nil {
		return nil
	}
	rf := reflect.ValueOf(val)
	if rf.Kind() != reflect.Map {
		return errors.New("value must be null or type map")
	}
	return nil
}

// reverseString reverses the input string character by character.
// Returns an empty string if the input is empty.
func reverseString(rawString string) string {
	if rawString == "" {
		return rawString
	}
	var reverseString string

	for i := len([]rune(rawString)); i > 0; i-- {
		var str = rawString[i-1]
		reverseString += string(str)
	}

	return reverseString
}

// ValidCitizenId validates a Thai citizen ID using the MOD 11 algorithm.
// The citizen ID must be exactly 13 digits long.
// Returns true if the ID is valid according to the checksum algorithm, false otherwise.
func ValidCitizenId(citizen string) bool {
	if len(citizen) != 13 {
		return false
	}

	var revString = reverseString(citizen)
	var total float64
	for index := 1; index < 13; index++ {
		var mul = index + 1
		var num, _ = strconv.Atoi(string([]rune(revString)[index]))
		var count = num * mul
		total = total + float64(count)
	}
	var mod = int(total) % 11
	var sub = 11 - mod
	var checkDigit = sub % 10

	var lastCitizen, _ = strconv.Atoi(string([]rune(revString)[0]))
	if lastCitizen == checkDigit {
		return true
	}

	return false
}

// IsCompany checks if the given citizen ID belongs to a company (juristic person).
// In Thailand, company IDs start with 0, while personal IDs start with 1-9.
// Returns true if the ID is valid and belongs to a company, false otherwise.
func IsCompany(citizen string) bool {
	if !ValidCitizenId(citizen) {
		return false
	}

	var num, err = strconv.Atoi(string([]rune(citizen)[0]))
	if err != nil {
		return false
	}
	if num == 0 {
		return true
	}

	return false
}
