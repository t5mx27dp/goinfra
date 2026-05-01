package logrusutil

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestJSON(
	t *testing.T,
	log func(logger *logrus.Logger),
	assert func(fields logrus.Fields),
) {
	var (
		out    bytes.Buffer
		fields logrus.Fields
	)

	logger := logrus.New()

	logger.Out = &out
	logger.Formatter = &logrus.JSONFormatter{}

	log(logger)

	err := json.Unmarshal(out.Bytes(), &fields)
	require.Nil(t, err)

	assert(fields)
}
