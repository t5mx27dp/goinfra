package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
		keys := []string{"key1", "key2"}

		manager := NewManager(keys[0])

		manager.WithKey(keys[1])

		require.Equal(t, keys[0], manager.Key())
		require.Equal(t, keys, manager.Keys())
	})

	t.Run("GenerateTraceID", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		require.Equal(t, 36, len(manager.GenerateTraceID()))

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		require.Equal(t, value, manager.GenerateTraceID())
	})

	t.Run("Trace_Context_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := manager.Trace(nil)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Trace_Context_Not_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = manager.Trace(ctx)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Trace_Context_Key_Exists", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = context.WithValue(ctx, key, "test")

		ctx = manager.Trace(ctx)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Parse_Context_Nil", func(t *testing.T) {
		key := "key"

		manager := NewManager(key)

		labels := manager.Parse(nil)

		require.Equal(t, 0, len(labels))
	})

	t.Run("Parse_Context_Not_Nil", func(t *testing.T) {
		key := "key"

		manager := NewManager(key)

		ctx := context.Background()

		labels := manager.Parse(ctx)

		require.Equal(t, "", labels[key])
	})

	t.Run("Parse_Context_Key_Exists", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		ctx := context.Background()

		ctx = context.WithValue(ctx, key, value)

		labels := manager.Parse(ctx)

		require.Equal(t, value, labels[key])
	})
}
