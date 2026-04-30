package tracing

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

var (
	globalManager      *Manager
	globalManagerMutex sync.RWMutex
)

func GetGlobalManager() *Manager {
	globalManagerMutex.RLock()
	defer globalManagerMutex.RUnlock()

	return globalManager
}

func SetGlobalManager(manager *Manager) {
	globalManagerMutex.Lock()
	defer globalManagerMutex.Unlock()

	globalManager = manager
}

type GenerateTraceID func() string

type contextKey string

type Manager struct {
	key  string
	keys []string

	generateTraceID GenerateTraceID

	mutex sync.RWMutex
}

func NewManager(key string) *Manager {
	return &Manager{
		key:  key,
		keys: []string{key},
	}
}

func (m *Manager) Key() string {
	return m.key
}

func (m *Manager) Keys() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.keys
}

func (m *Manager) WithKey(key string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.keys = append(m.keys, key)
}

func (m *Manager) SetTraceIDGenerator(fn GenerateTraceID) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.generateTraceID = fn
}

func (m *Manager) GenerateTraceID() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if m.generateTraceID != nil {
		return m.generateTraceID()
	}
	return uuid.New().String()
}

func (m *Manager) Trace(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	return context.WithValue(ctx, contextKey(m.key), m.GenerateTraceID())
}

func (m *Manager) TraceWithValue(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

func (m *Manager) Parse(ctx context.Context) map[string]string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if ctx == nil {
		return nil
	}

	labels := make(map[string]string)

	for _, key := range m.keys {
		value := ctx.Value(contextKey(key))
		if value == nil {
			labels[key] = ""
		} else {
			labels[key], _ = value.(string)
		}
	}

	return labels
}
