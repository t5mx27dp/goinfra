package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
)

func GinMiddleware(manager *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := manager.Trace(c.Request.Context())

		for _, key := range manager.Keys() {
			if key == manager.Key() {
				continue
			}
			ctx = context.WithValue(ctx, key, c.GetHeader(key))
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
