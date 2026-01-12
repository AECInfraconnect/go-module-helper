# Go Module Helper

A comprehensive Go helper library for common utility functions, validation, context management, and API responses.

## Installation

```bash
go get github.com/AECInfraconnect/go-module-helper
```

## Features

### Environment Variables

- `GetENV(key, defaultValue string) string` - Get environment variable with fallback

### Array/Slice Utilities

- `InArray[T](needle T, haystack []T) bool` - Check if value exists in slice
- `Contains[T](haystack []T, needle T) bool` - Alias for InArray

### String Utilities

- `IsEmpty(s string) bool` - Check if string is empty or whitespace
- `Truncate(s string, maxLength int, suffix ...string) string` - Truncate string with optional suffix

### File/MIME Type

- `GetMimeType(filename string) string` - Get MIME type from filename

### Type Conversion

- `ValidateTypeString(value interface{}) (string, bool)` - Convert to string
- `ValidateTypeInt(value interface{}) (int, bool)` - Convert to int

### Pointer Helpers

- `Ptr[T](v T) *T` - Generic pointer creator
- `PtrString(s string) *string` - String pointer
- `PtrInt(i int) *int` - Int pointer
- `PtrBool(b bool) *bool` - Bool pointer
- `ValueOr[T](ptr *T, defaultValue T) T` - Safe pointer dereference
- `StringValue(*string) string` - Get string value or empty
- `IntValue(*int) int` - Get int value or 0
- `BoolValue(*bool) bool` - Get bool value or false

### Error Handling

- `Must[T](value T, err error) T` - Panic on error
- `Must2[T1, T2](v1 T1, v2 T2, err error) (T1, T2)` - Must for 2 return values
- `IgnoreError[T](value T, err error) T` - Return value, ignore error
- `IgnoreValue[T](value T, err error) error` - Return error, ignore value

### Validation Functions

- `ValidateTypeString(val interface{}) error` - Validate string type
- `ValidateTypeInt(val interface{}) error` - Validate int type
- `ValidateTypeFloat(val interface{}) error` - Validate float type
- `ValidateTypeBool(val interface{}) error` - Validate bool type
- `ValidateTypeMap(val interface{}) error` - Validate map type
- `ValidateTypeSlice(val interface{}) error` - Validate slice type
- `ValidateTypeUUID(val interface{}) error` - Validate UUID format
- `ValidateUUIDOrIDZero(val interface{}) error` - Validate UUID or "0"
- `ValidateOnlyThaiLetterNumeric(val interface{}) error` - Thai letters + numbers only
- `ValidCitizenId(citizen string) bool` - Validate Thai citizen ID
- `IsCompany(citizen string) bool` - Check if citizen ID is company

### Gin Context Helpers

- `GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool)` - Get user UUID
- `GetRequestIDFromContext(c *gin.Context) string` - Get request ID
- `GetIPAddress(c *gin.Context) string` - Get client IP
- `GetUserAgent(c *gin.Context) string` - Get user agent
- `IsAPIKeyAuth(c *gin.Context) bool` - Check API key authentication

### API Response Helpers

- `SuccessResponse(c *gin.Context, statusCode int, data interface{})` - Send success response
- `ErrorResponse(c *gin.Context, statusCode int, code, message string)` - Send error response
- `ValidationErrorResponse(c *gin.Context, err error)` - Send validation error

### Timestamp Utilities

- `Timestamp` type with custom marshaling/unmarshaling
- `NewTimestampFromString(s string) Timestamp`
- `NewTimestampFromTime(t time.Time) Timestamp`
- `NewTimestampAddDayFromTime(t time.Time, years, months, days int) Timestamp`

## Usage Examples

```go
import "github.com/AECInfraconnect/go-module-helper/helper"

// Environment
port := helper.GetENV("PORT", "8080")

// Array operations
exists := helper.Contains([]string{"a", "b", "c"}, "b") // true

// String utilities
helper.IsEmpty("  ") // true
helper.Truncate("Hello World", 5) // "He..."

// Pointers
name := helper.PtrString("John")
age := helper.ValueOr(nil, 18) // returns 18

// Error handling
file := helper.Must(os.Open("file.txt"))

// Validation
err := helper.ValidateTypeUUID("550e8400-e29b-41d4-a716-446655440000")
valid := helper.ValidCitizenId("1234567890123")

// Gin response
helper.SuccessResponse(c, 200, gin.H{"message": "OK"})
helper.ErrorResponse(c, 400, "INVALID_INPUT", "Invalid data")
```

## Dependencies

- `github.com/gin-gonic/gin` - Web framework
- `github.com/google/uuid` - UUID support
