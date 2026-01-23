package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys
const (
	ContextKeyUserID     = "user_id"
	ContextKeyAPIKeyAuth = "apiKey"
)

// GetUserIDFromContext retrieves user ID from context
func GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return uuid.Nil, false
	}

	id, ok := userID.(uuid.UUID)
	return id, ok
}

// GetIPAddress retrieves client IP address
func GetIPAddress(c *gin.Context) string {
	return c.ClientIP()
}

// GetUserAgent retrieves user agent
func GetUserAgent(c *gin.Context) string {
	return c.Request.UserAgent()
}

// IsAPIKeyAuth checks if request was authenticated via API Key
func IsAPIKeyAuth(c *gin.Context) bool {
	apiKeyAuth, exists := c.Get(ContextKeyAPIKeyAuth)
	if !exists {
		return false
	}
	isAuth, ok := apiKeyAuth.(bool)
	return ok && isAuth
}
