package tracing

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogrusHook(t *testing.T) {
	key := "key"
	value := "value"

	manager := NewManager(key)

	manager.SetTraceIDGenerator(func() string {
		return value
	})

	var (
		buffer bytes.Buffer
		fields logrus.Fields
	)

	log := logrus.New()

	log.SetOutput(&buffer)
	log.SetFormatter(&logrus.JSONFormatter{})

	log.AddHook(NewLogrusHook(manager))

	ctx := manager.Trace(nil)

	log.WithContext(ctx).Error("err")

	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	require.Equal(t, value, fields[key])
}
