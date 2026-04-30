package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/stretchr/testify/require"
)

func TestGinMiddleware(t *testing.T) {
	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}

	manager := NewManager(keys[0])

	manager.WithKey(keys[1])

	manager.SetTraceIDGenerator(func() string {
		return values[0]
	})

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(GinMiddleware(manager))

	router.GET("/", func(c *gin.Context) {
		ctx := c.Request.Context()

		labels := manager.Parse(ctx)

		require.Equal(t, values[0], labels[keys[0]])
		require.Equal(t, values[1], labels[keys[1]])
	})

	server := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	go func() {
		options := &grequests.RequestOptions{
			Headers: map[string]string{
				keys[1]: values[1],
			},
		}
		_, err := grequests.NewSession(options).Get("http://127.0.0.1:8000/", nil)
		require.Nil(t, err)

		err = server.Shutdown(context.Background())
		require.Nil(t, err)
	}()

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		t.Error(err)
		return
	}
}
