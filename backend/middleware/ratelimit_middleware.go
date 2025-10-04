package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPClients stores a map of IP addresses to their respective rate limiters.
var IPClients = make(map[string]*rate.Limiter)
var mu sync.Mutex

// RateLimitMiddleware returns a Gin middleware that limits requests per IP address.
func RateLimitMiddleware(limit float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		if _, found := IPClients[ip]; !found {
			IPClients[ip] = rate.NewLimiter(rate.Limit(limit), burst)
		}
		limiter := IPClients[ip]
		mu.Unlock()

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
		c.Next()
	}
}
