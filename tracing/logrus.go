package tracing

import (
	"github.com/sirupsen/logrus"
)

type LogrusHook struct {
	m *Manager
}

func NewLogrusHook(m *Manager) logrus.Hook {
	return &LogrusHook{
		m: m,
	}
}

func (h *LogrusHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *LogrusHook) Fire(entry *logrus.Entry) error {
	for key, value := range h.m.Parse(entry.Context) {
		entry.Data[key] = value
	}
	return nil
}
