package middleware

import (
	"net/http"
	"strings"

	"github.com/AECInfraconnect/go-module-helper/helper"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// APIKeyAuthMiddleware validates API keys for service-to-service authentication.
//
// Reads the API-Key header and compares it against the provided apiKey.
// Sets helper.ContextKeyAPIKeyAuth to true on successful authentication.
//
// Example:
//
//	service := r.Group("/service")
//	service.Use(middleware.APIKeyAuthMiddleware("secret-key"))
func APIKeyAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if API Key is not configured
		if apiKey == "" {
			c.Next()
			return
		}

		apiKeyHeader := c.GetHeader("API-Key")
		if apiKeyHeader == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "MISSING_API_KEY", "API-Key header is required")
			c.Abort()
			return
		}

		// Validate API Key
		if apiKeyHeader != apiKey {
			helper.ErrorResponse(c, http.StatusUnauthorized, "INVALID_API_KEY", "Invalid API Key")
			c.Abort()
			return
		}

		// Set a flag indicating API Key authentication was successful
		c.Set(helper.ContextKeyAPIKeyAuth, true)
		c.Next()
	}
}

// APIKeyOrJWTAuthMiddleware accepts either API Key or JWT token authentication.
//
// Attempts API key authentication first, then falls back to JWT.
// Useful for endpoints that need to support both service and user authentication.
//
// Example:
//
//	r.Use(middleware.APIKeyOrJWTAuthMiddleware("service-key", "jwt-secret"))
func APIKeyOrJWTAuthMiddleware(apiKey string, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try API Key first
		apiKeyHeader := c.GetHeader("API-Key")
		if apiKeyHeader != "" && apiKey != "" && apiKeyHeader == apiKey {
			c.Set(helper.ContextKeyAPIKeyAuth, true)
			c.Next()
			return
		}

		// Try JWT if API Key not provided or invalid
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "MISSING_AUTH", "API-Key or Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			helper.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Token must be in Bearer format")
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			helper.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Token is invalid or expired")
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Extract user ID
			if userIDStr, ok := claims["user_id"].(string); ok {
				if userID, err := uuid.Parse(userIDStr); err == nil {
					c.Set(helper.ContextKeyUserID, userID)
				}
			}
		}

		c.Next()
	}
}

// JWTAuthMiddleware validates JWT tokens for user authentication.
//
// Expects "Authorization: Bearer <token>" header format.
// Extracts user_id from JWT claims and stores it in context.
//
// Example:
//
//	auth := r.Group("/auth")
//	auth.Use(middleware.JWTAuthMiddleware("jwt-secret"))
func JWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			helper.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Token must be in Bearer format")
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			helper.ErrorResponse(c, http.StatusUnauthorized, "INVALID_TOKEN", "Token is invalid or expired")
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Extract user ID
			if userIDStr, ok := claims["user_id"].(string); ok {
				if userID, err := uuid.Parse(userIDStr); err == nil {
					c.Set(helper.ContextKeyUserID, userID)
				}
			}
		}

		c.Next()
	}
}

// OptionalJWTAuthMiddleware validates JWT tokens if present but doesn't require them.
//
// Allows both authenticated and anonymous requests.
// Handlers can check context for user_id to determine authentication status.
//
// Example:
//
//	r.Use(middleware.OptionalJWTAuthMiddleware("jwt-secret"))
func OptionalJWTAuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if userIDStr, ok := claims["user_id"].(string); ok {
					if userID, err := uuid.Parse(userIDStr); err == nil {
						c.Set(helper.ContextKeyUserID, userID)
					}
				}
			}
		}

		c.Next()
	}
}

// RequirePermission checks if the authenticated user has a specific permission.
//
// This is a placeholder for custom permission logic.
// Implement by querying your database or permission service.
func RequirePermission(scope string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement permission check logic
		// This would typically query the database to check user permissions
		c.Next()
	}
}

// RequireRole checks if the authenticated user has a specific role.
//
// This is a placeholder for custom role logic.
// Implement by querying your database or role service.
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement role check logic
		// This would typically query the database to check user roles
		c.Next()
	}
}
