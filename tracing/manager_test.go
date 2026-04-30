package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGlobalManager(t *testing.T) {
	t.Cleanup(func() {
		SetGlobalManager(nil)
	})

	require.Nil(t, globalManager)

	manager := NewManager("")

	SetGlobalManager(manager)

	require.Equal(t, globalManager, GetGlobalManager())
}

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

	t.Run("Trace_Parse_Context_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := manager.Trace(nil)

		labels := manager.Parse(ctx)

		require.Equal(t, value, labels[key])
	})

	t.Run("Trace_Parse_Context_Not_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = manager.Trace(ctx)

		labels := manager.Parse(ctx)

		require.Equal(t, value, labels[key])
	})

	t.Run("Trace_Parse_Context_Key_Exists", func(t *testing.T) {
		key := "key"
		value := "value"

		manager := NewManager(key)

		manager.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = manager.TraceWithValue(ctx, key, "test")

		ctx = manager.Trace(ctx)

		labels := manager.Parse(ctx)

		require.Equal(t, value, labels[key])
	})
}
