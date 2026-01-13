# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added

#### Convert Package (v0.2.0)

- **Type Conversion** (`convert/type.go`)
  - `ToString()` - Convert any value to string representation
  - `ToInt()` - Convert various types to int with error handling
  - `ToInt64()` - Convert various types to 64-bit integer
  - `ToFloat64()` - Convert various types to 64-bit float
  - `ToBool()` - Flexible boolean conversion from multiple types
  - `ToStringSlice()` - Convert values to string slice
  - `ToIntSlice()` - Convert values to integer slice with validation

- **JSON Operations** (`convert/json.go`)
  - `ToJSON()` - Marshal any value to JSON string
  - `ToJSONIndent()` - Marshal with custom indentation
  - `FromJSON()` - Unmarshal JSON string to interface{}
  - `FromJSONTo()` - Unmarshal JSON to specific type
  - `StructToMap()` - Convert struct to map[string]interface{}
  - `MapToStruct()` - Convert map to struct
  - `ToJSONBytes()` - Efficient JSON byte conversion
  - `FromJSONBytes()` - Parse JSON from bytes
  - `IsValidJSON()` - Validate JSON string format

- **Enum Handling** (`convert/enum.go`)
  - `EnumToString()` - Convert enum to string
  - `StringToEnum()` - Convert string to enum with mapping
  - `EnumToInt()` - Convert enum to integer
  - `IntToEnum()` - Convert integer to enum with mapping
  - `EnumExists()` - Check if enum value exists (string)
  - `EnumExistsInt()` - Check if enum value exists (int)
  - `GetEnumKeys()` - Extract all keys from enum map
  - `GetEnumValues()` - Extract all values from enum map
  - `NormalizeEnumString()` - Normalize string for comparison
  - `EnumMatch()` - Case-insensitive enum matching
  - `ValidateEnum()` - Validate string enum with error messages
  - `ValidateEnumInt()` - Validate integer enum with error messages

- Comprehensive GoDoc comments with examples for all convert functions

## [0.1.0] - 2025-01-XX

### Added

#### Helper Package

- **Environment** (`helper/env.go`)
  - `GetEnv()` - Get environment variable with default value support

- **Context Helpers** (`helper/context.go`)
  - `GetUserID()` - Extract user ID from Gin context
  - `GetUserRole()` - Extract user role from Gin context
  - `GetAuthToken()` - Extract authentication token from Gin context
  - `SetUserContext()` - Set user information in Gin context
  - `GetContextValue()` - Generic context value retrieval

- **Response Helpers** (`helper/response.go`)
  - `SuccessResponse()` - Standard success response structure
  - `ErrorResponse()` - Standard error response structure
  - `PaginatedResponse()` - Paginated data response
  - `ValidationErrorResponse()` - Validation error formatting
  - Response helper functions for Gin framework

- **Validation** (`helper/validate.go`)
  - `ValidateEmail()` - Email format validation
  - `ValidatePhone()` - Phone number validation (Thai format)
  - `ValidateURL()` - URL format validation
  - `ValidatePassword()` - Password strength validation
  - `ValidateRequired()` - Required field validation
  - `ValidateMinLength()` - Minimum length validation
  - `ValidateMaxLength()` - Maximum length validation
  - `ValidateRange()` - Numeric range validation
  - `ValidateAlphanumeric()` - Alphanumeric validation
  - Comprehensive validation helper functions

- **Pagination** (`helper/pagination.go`)
  - `Paginate()` - GORM pagination helper
  - `GetPaginationParams()` - Extract pagination parameters from request
  - `CalculatePaginationMeta()` - Calculate pagination metadata
  - Pagination helper with GORM support

- **Array Utilities** (`helper/array.go`)
  - `InArray()` - Check if value exists in array with position
  - `Contains()` - Check if string exists in string slice

- **MIME Type Utilities** (`helper/mime.go`)
  - `GetMIMEType()` - Get MIME type from file extension
  - `GetMIMECategory()` - Get MIME category (image, video, document, etc.)
  - `IsImageMIME()` - Check if MIME type is an image
  - `IsVideoMIME()` - Check if MIME type is a video
  - `IsAudioMIME()` - Check if MIME type is an audio
  - `IsDocumentMIME()` - Check if MIME type is a document
  - MIME type utilities with document type detection

- **Timestamp Utilities** (`helper/timestamp.go`)
  - `CurrentTimestamp()` - Get current Unix timestamp
  - `FormatTimestamp()` - Format timestamp to string
  - `ParseTimestamp()` - Parse string to timestamp

### Documentation

- Comprehensive README.md with installation guide, usage examples, and API documentation
- Package-level documentation for all modules
- Function-level GoDoc comments with practical examples

### Infrastructure

- Initial project setup
- Go module configuration
- Git repository initialization

## Release Notes

### v0.2.0 - Convert Package

This release introduces a comprehensive type conversion package with support for:

- Flexible type conversions between primitive types
- JSON serialization/deserialization utilities
- Enum validation and conversion helpers
- Full documentation with examples

### v0.1.0 - Initial Release

First stable release with helper utilities for:

- Gin framework integration (context, response)
- Validation helpers
- Pagination support
- Array and MIME type utilities
- Environment variable management

---

## Migration Guide

### From No Version to v0.1.0

```go
// Install the module
go get github.com/reemwsw/go-module-helper

// Import helpers
import (
    "github.com/reemwsw/go-module-helper/helper"
)
```

### From v0.1.0 to v0.2.0

```go
// Import the new convert package
import (
    "github.com/reemwsw/go-module-helper/convert"
    "github.com/reemwsw/go-module-helper/helper"
)

// Use type conversion
str := convert.ToString(123)
num, err := convert.ToInt("456")

// Use JSON utilities
jsonStr, err := convert.ToJSON(myStruct)

// Use enum validation
status, err := convert.ValidateEnum("ACTIVE", validStatuses)
```

---

## Contributors

- [@reemwsw](https://github.com/reemwsw) - Creator and maintainer

## License

This project follows semantic versioning and is maintained as an open-source utility library.
