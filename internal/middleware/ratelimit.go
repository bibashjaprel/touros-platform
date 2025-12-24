package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter *rate.Limiter
}

var rateLimiters = make(map[string]*rateLimiter)

func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter, exists := rateLimiters[ip]
		if !exists {
			limiter = &rateLimiter{
				limiter: rate.NewLimiter(rate.Limit(rps), burst),
			}
			rateLimiters[ip] = limiter
		}

		if !limiter.limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

