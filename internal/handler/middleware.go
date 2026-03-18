package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/VAGRAMCHIC/vless_reality_agent/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func APIKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("X-API-KEY") != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func RequestIDMiddleware(log *utils.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, utils.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", requestID)

		start := time.Now()
		c.Next()
		duration := time.Since(start)

		log.Info(ctx, "http_request", map[string]interface{}{
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": duration.Milliseconds(),
			"ip":       c.ClientIP(),
		})
	}
}
