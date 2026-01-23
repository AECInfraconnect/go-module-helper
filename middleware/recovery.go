package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/AECInfraconnect/go-module-helper/helper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware recovers from panics and logs the error
func RecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get request ID for tracing
				requestID := GetRequestID(c)

				// Get stack trace
				stack := string(debug.Stack())

				// Log the panic
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      err,
					"stack":      stack,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"client_ip":  c.ClientIP(),
				}).Error("Panic recovered")

				// Return error response
				helper.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()

		c.Next()
	}
}

// SimpleRecoveryMiddleware recovers from panics without logging
func SimpleRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)

				// Simple console output
				fmt.Printf("[PANIC] Request ID: %s | Error: %v\n", requestID, err)
				fmt.Printf("Stack trace:\n%s\n", debug.Stack())

				helper.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RecoveryWithCustomHandlerMiddleware allows custom panic handling
func RecoveryWithCustomHandlerMiddleware(handler func(*gin.Context, interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				handler(c, err)
				c.Abort()
			}
		}()

		c.Next()
	}
}

// RecoveryWithNotificationMiddleware recovers and sends notification
// You can integrate with services like Sentry, Slack, etc.
func RecoveryWithNotificationMiddleware(logger *logrus.Logger, notifyFunc func(requestID string, err interface{}, stack string)) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(c)
				stack := string(debug.Stack())

				// Log the panic
				logger.WithFields(logrus.Fields{
					"request_id": requestID,
					"error":      err,
					"stack":      stack,
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"client_ip":  c.ClientIP(),
				}).Error("Panic recovered")

				// Send notification
				if notifyFunc != nil {
					go notifyFunc(requestID, err, stack)
				}

				helper.ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "An unexpected error occurred")
				c.Abort()
			}
		}()

		c.Next()
	}
}
