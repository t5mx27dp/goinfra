package tracing

import (
	"context"

	"github.com/gin-gonic/gin"
)

func GinMiddleware(m *Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := m.Trace(c.Request.Context())

		for _, key := range m.Keys() {
			if key == m.Key() {
				continue
			}
			ctx = context.WithValue(ctx, key, c.GetHeader(key))
		}

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
