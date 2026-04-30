package tracing

import (
	"github.com/gin-gonic/gin"
)

func GinMiddleware(manager *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		if c.GetHeader(manager.Key()) == "" {
			ctx = manager.Trace(ctx)
		}

		for _, key := range manager.Keys() {
			if key == manager.Key() {
				continue
			}
			ctx = manager.TraceWithValue(ctx, key, c.GetHeader(key))
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
