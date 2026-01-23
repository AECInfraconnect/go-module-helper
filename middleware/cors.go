// Package middleware provides production-ready Gin middlewares for authentication,
// rate limiting, CORS handling, request tracking, and form parsing.
//
// Basic usage:
//
//	r := gin.Default()
//	r.Use(middleware.CORSMiddleware())
//	r.Use(middleware.RequestIDMiddleware())
//
// All middlewares are designed to be composable and work together seamlessly.
package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles Cross-Origin Resource Sharing (CORS).
//
// Allows all origins, credentials, and common HTTP methods.
// Automatically handles OPTIONS preflight requests.
//
// Example:
//
//	r.Use(middleware.CORSMiddleware())
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
