package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/AECInfraconnect/go-module-helper/helper"
	"github.com/gin-gonic/gin"
)

// RateLimiter implements a token bucket rate limiting algorithm.
//
// Tokens are refilled at a constant rate over time.
// Each request consumes one token.
type RateLimiter struct {
	tokens         int
	maxTokens      int
	refillRate     time.Duration
	lastRefillTime time.Time
	mu             sync.Mutex
}

// RateLimiterStore manages rate limiters for multiple clients.
//
// Automatically cleans up inactive limiters every 10 minutes.
type RateLimiterStore struct {
	limiters map[string]*RateLimiter
	mu       sync.RWMutex
}

// NewRateLimiterStore creates a new rate limiter store with automatic cleanup.
func NewRateLimiterStore() *RateLimiterStore {
	store := &RateLimiterStore{
		limiters: make(map[string]*RateLimiter),
	}
	// Start cleanup goroutine to remove inactive limiters
	go store.cleanup()
	return store
}

// cleanup removes inactive limiters every 10 minutes
func (s *RateLimiterStore) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for key, limiter := range s.limiters {
			limiter.mu.Lock()
			// Remove limiters that haven't been used in 30 minutes
			if now.Sub(limiter.lastRefillTime) > 30*time.Minute {
				delete(s.limiters, key)
			}
			limiter.mu.Unlock()
		}
		s.mu.Unlock()
	}
}

// GetLimiter retrieves or creates a rate limiter for the specified client.
func (s *RateLimiterStore) GetLimiter(clientID string, maxTokens int, refillRate time.Duration) *RateLimiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	if limiter, exists := s.limiters[clientID]; exists {
		return limiter
	}

	limiter := &RateLimiter{
		tokens:         maxTokens,
		maxTokens:      maxTokens,
		refillRate:     refillRate,
		lastRefillTime: time.Now(),
	}
	s.limiters[clientID] = limiter
	return limiter
}

// Allow checks if a request is allowed under the current rate limit.
//
// Returns true if a token is available and consumes it.
// Returns false if no tokens are available.
func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastRefillTime)

	// Refill tokens based on elapsed time
	tokensToAdd := int(elapsed / r.refillRate)
	if tokensToAdd > 0 {
		r.tokens += tokensToAdd
		if r.tokens > r.maxTokens {
			r.tokens = r.maxTokens
		}
		r.lastRefillTime = now
	}

	// Check if we have tokens available
	if r.tokens > 0 {
		r.tokens--
		return true
	}

	return false
}

// RateLimitMiddleware creates a rate limiting middleware based on user ID or IP.
//
// Limits requests per authenticated user (if JWT middleware is used) or per IP address.
//
// Example:
//
//	// 100 requests per minute
//	r.Use(middleware.RateLimitMiddleware(100, time.Minute))
func RateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	store := NewRateLimiterStore()
	refillRate := window / time.Duration(maxRequests)

	return func(c *gin.Context) {
		// Get client identifier (IP address or user ID if authenticated)
		clientID := c.ClientIP()

		// If user is authenticated, use user ID instead
		if userID, exists := c.Get(helper.ContextKeyUserID); exists {
			clientID = userID.(string)
		}

		// Get or create limiter for this client
		limiter := store.GetLimiter(clientID, maxRequests, refillRate)

		// Check if request is allowed
		if !limiter.Allow() {
			helper.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimitMiddleware creates a rate limiting middleware based solely on IP address.
//
// Example:
//
//	// 1000 requests per hour per IP
//	r.Use(middleware.IPRateLimitMiddleware(1000, time.Hour))
func IPRateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	store := NewRateLimiterStore()
	refillRate := window / time.Duration(maxRequests)

	return func(c *gin.Context) {
		clientID := c.ClientIP()
		limiter := store.GetLimiter(clientID, maxRequests, refillRate)

		if !limiter.Allow() {
			helper.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Too many requests from this IP. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// APIKeyRateLimitMiddleware creates a rate limiting middleware based on API keys.
//
// Reads the API-Key header for client identification.
// Falls back to IP address if no API key is provided.
//
// Example:
//
//	// 10000 requests per day per API key
//	r.Use(middleware.APIKeyRateLimitMiddleware(10000, 24*time.Hour))
func APIKeyRateLimitMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	store := NewRateLimiterStore()
	refillRate := window / time.Duration(maxRequests)

	return func(c *gin.Context) {
		apiKey := c.GetHeader("API-Key")
		if apiKey == "" {
			apiKey = c.ClientIP() // Fallback to IP if no API key
		}

		limiter := store.GetLimiter(apiKey, maxRequests, refillRate)

		if !limiter.Allow() {
			helper.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "API rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}

// CustomRateLimitMiddleware creates a rate limiting middleware with custom client ID extraction.
//
// The getClientID function determines how to identify clients.
// Falls back to IP address if the function returns an empty string.
//
// Example:
//
//	getID := func(c *gin.Context) string {
//	    return c.GetHeader("X-Tenant-ID")
//	}
//	r.Use(middleware.CustomRateLimitMiddleware(500, time.Minute, getID))
func CustomRateLimitMiddleware(maxRequests int, window time.Duration, getClientID func(*gin.Context) string) gin.HandlerFunc {
	store := NewRateLimiterStore()
	refillRate := window / time.Duration(maxRequests)

	return func(c *gin.Context) {
		clientID := getClientID(c)
		if clientID == "" {
			clientID = c.ClientIP() // Fallback to IP
		}

		limiter := store.GetLimiter(clientID, maxRequests, refillRate)

		if !limiter.Allow() {
			helper.ErrorResponse(c, http.StatusTooManyRequests, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded. Please try again later.")
			c.Abort()
			return
		}

		c.Next()
	}
}
