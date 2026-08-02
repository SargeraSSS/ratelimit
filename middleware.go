package ratelimit

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(limiter *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		allowed, retryAfter := limiter.Allow(clientIP)
		if !allowed {
			c.Header("Retry-After", retryAfterSeconds(retryAfter))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}

// retryAfterSeconds renders a wait as the whole seconds the Retry-After header
// requires. It rounds up: rounding down would invite the client back before the
// tokens exist, and the retry would be refused a second time.
func retryAfterSeconds(wait time.Duration) string {
	seconds := int64(math.Ceil(wait.Seconds()))
	if seconds < 1 {
		seconds = 1
	}
	return strconv.FormatInt(seconds, 10)
}
