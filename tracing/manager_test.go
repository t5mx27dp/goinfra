package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	t.Run("Keys", func(t *testing.T) {
		keys := []string{"key1", "key2"}

		m := NewManager(keys[0])

		m.WithKey(keys[1])

		require.Equal(t, keys[0], m.Key())
		require.Equal(t, keys, m.Keys())
	})

	t.Run("GenerateTraceID", func(t *testing.T) {
		key := "key"
		value := "value"

		m := NewManager(key)

		require.Equal(t, 36, len(m.GenerateTraceID()))

		m.SetTraceIDGenerator(func() string {
			return value
		})

		require.Equal(t, value, m.GenerateTraceID())
	})

	t.Run("Trace_Context_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		m := NewManager(key)

		m.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := m.Trace(nil)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Trace_Context_Not_Nil", func(t *testing.T) {
		key := "key"
		value := "value"

		m := NewManager(key)

		m.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = m.Trace(ctx)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Trace_Context_Key_Exists", func(t *testing.T) {
		key := "key"
		value := "value"

		m := NewManager(key)

		m.SetTraceIDGenerator(func() string {
			return value
		})

		ctx := context.Background()

		ctx = context.WithValue(ctx, key, "test")

		ctx = m.Trace(ctx)

		require.Equal(t, value, ctx.Value(key))
	})

	t.Run("Parse_Context_Nil", func(t *testing.T) {
		key := "key"

		m := NewManager(key)

		labels := m.Parse(nil)

		require.Equal(t, 0, len(labels))
	})

	t.Run("Parse_Context_Not_Nil", func(t *testing.T) {
		key := "key"

		m := NewManager(key)

		ctx := context.Background()

		labels := m.Parse(ctx)

		require.Equal(t, "", labels[key])
	})

	t.Run("Parse_Context_Key_Exists", func(t *testing.T) {
		key := "key"
		value := "value"

		m := NewManager(key)

		ctx := context.Background()

		ctx = context.WithValue(ctx, key, value)

		labels := m.Parse(ctx)

		require.Equal(t, value, labels[key])
	})
}
