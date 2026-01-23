package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// LoggerMiddleware logs HTTP requests with request ID
func LoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Get request ID
		requestID := GetRequestID(c)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Create log entry
		entry := logger.WithFields(logrus.Fields{
			"request_id":  requestID,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"query":       c.Request.URL.RawQuery,
			"status":      statusCode,
			"latency":     latency.String(),
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
			"error":       c.Errors.ByType(gin.ErrorTypePrivate).String(),
		})

		// Add user ID if authenticated
		if userID, exists := c.Get("user_id"); exists {
			entry = entry.WithField("user_id", userID)
		}

		// Log based on status code
		msg := fmt.Sprintf("%s %s - %d", c.Request.Method, c.Request.URL.Path, statusCode)

		switch {
		case statusCode >= 500:
			entry.Error(msg)
		case statusCode >= 400:
			entry.Warn(msg)
		case statusCode >= 300:
			entry.Info(msg)
		default:
			entry.Info(msg)
		}
	}
}

// SimpleLoggerMiddleware is a lightweight logger
func SimpleLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := GetRequestID(c)

		c.Next()

		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		// Simple colored console output
		fmt.Printf("[%s] %s | %3d | %13v | %s | %s %s\n",
			time.Now().Format("2006-01-02 15:04:05"),
			requestID,
			statusCode,
			latency,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)
	}
}

// AccessLoggerMiddleware logs only successful requests
func AccessLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := GetRequestID(c)

		c.Next()

		// Only log successful requests (2xx and 3xx)
		statusCode := c.Writer.Status()
		if statusCode < 400 {
			latency := time.Since(startTime)

			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     statusCode,
				"latency_ms": latency.Milliseconds(),
				"client_ip":  c.ClientIP(),
			}).Info("Request processed")
		}
	}
}

// ErrorLoggerMiddleware logs only errors (4xx and 5xx)
func ErrorLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)

		c.Next()

		// Only log error requests
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			logger.WithFields(logrus.Fields{
				"request_id": requestID,
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"status":     statusCode,
				"client_ip":  c.ClientIP(),
				"error":      c.Errors.ByType(gin.ErrorTypePrivate).String(),
			}).Error("Request failed")
		}
	}
}

// DetailedLoggerMiddleware logs request and response details
func DetailedLoggerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := GetRequestID(c)

		// Log request
		logger.WithFields(logrus.Fields{
			"request_id":   requestID,
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"client_ip":    c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"content_type": c.ContentType(),
		}).Info("Incoming request")

		c.Next()

		// Log response
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()

		entry := logger.WithFields(logrus.Fields{
			"request_id":  requestID,
			"status":      statusCode,
			"latency":     latency.String(),
			"latency_ms":  latency.Milliseconds(),
			"body_size":   c.Writer.Size(),
		})

		if statusCode >= 400 {
			entry.Error("Request completed with error")
		} else {
			entry.Info("Request completed")
		}
	}
}
