package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	// RequestIDHeader is the HTTP header key for request ID.
	RequestIDHeader = "X-Request-ID"

	// RequestIDKey is the Gin context key for storing request ID.
	RequestIDKey = "request_id"
)

// RequestIDMiddleware generates or extracts request IDs for distributed tracing.
//
// If the client provides an X-Request-ID header, it will be used.
// Otherwise, a new UUID v4 will be generated.
// The request ID is stored in both the Gin context and response header.
//
// Example:
//
//	r.Use(middleware.RequestIDMiddleware())
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get request ID from header
		requestID := c.GetHeader(RequestIDHeader)

		// If not provided, generate a new one
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set(RequestIDKey, requestID)

		// Set request ID in response header for client tracking
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}

// ForceNewRequestIDMiddleware always generates a new request ID.
//
// Use this when you don't want to trust client-provided request IDs,
// such as for financial transactions or security-critical operations.
//
// Example:
//
//	webhook.Use(middleware.ForceNewRequestIDMiddleware())
func ForceNewRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Always generate new request ID
		requestID := uuid.New().String()

		// Set in context and response header
		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID retrieves the request ID from Gin context.
//
// Returns empty string if request ID middleware was not applied.
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}

// RequestIDWithPrefixMiddleware generates request IDs with a custom prefix.
//
// Example with prefix "api":
//
//	"api-550e8400-e29b-41d4-a716-446655440000"
//
// Usage:
//
//	r.Use(middleware.RequestIDWithPrefixMiddleware("api"))
func RequestIDWithPrefixMiddleware(prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)

		if requestID == "" {
			if prefix != "" {
				requestID = prefix + "-" + uuid.New().String()
			} else {
				requestID = uuid.New().String()
			}
		}

		c.Set(RequestIDKey, requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}
