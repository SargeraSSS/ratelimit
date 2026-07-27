package ratelimit

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRateLimitMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)
	limiter := NewLimiter(2, 1)
	defer limiter.Close()

	router := gin.New()
	router.Use(RateLimitMiddleware(limiter))
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	doRequest := func() int {
		req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)
		return recorder.Code
	}
	if status := doRequest(); status != http.StatusOK {
		t.Errorf("request 1: want %d, got %d", http.StatusOK, status)
	}
	if status := doRequest(); status != http.StatusOK {
		t.Errorf("request 2: want %d, got %d", http.StatusOK, status)
	}
	if status := doRequest(); status != http.StatusTooManyRequests {
		t.Errorf("request 3: want %d, got %d", http.StatusTooManyRequests, status)
	}
}
