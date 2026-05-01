package logrushook

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/t5mx27dp/goinfra/logrusutil"
)

func TestReportCallerHook(t *testing.T) {
	t.Run("Levels", func(t *testing.T) {
		hook := NewReportCallerHook("")

		hook.SetLevels(logrus.AllLevels)

		require.Equal(t, logrus.AllLevels, hook.Levels())
	})

	t.Run("Location", func(t *testing.T) {
		logrusutil.TestJSON(
			t,
			func(logger *logrus.Logger) {
				rootPath, err := os.Getwd()
				require.Nil(t, err)

				hook := NewReportCallerHook(rootPath)

				logger.AddHook(hook)

				logger.Error("error")
			},
			func(fields logrus.Fields) {
				require.Equal(t, "reportcaller_test.go:35", fields["Location"])
			},
		)
	})

	t.Run("Location_SetLocationBuilder", func(t *testing.T) {
		logrusutil.TestJSON(
			t,
			func(logger *logrus.Logger) {
				rootPath, err := os.Getwd()
				require.Nil(t, err)

				hook := NewReportCallerHook(rootPath)

				hook.SetLocationBuilder(func(filePath string, line int) string {
					return strings.Replace(filePath, rootPath+"/", "", -1) + "#" + strconv.Itoa(line)
				})

				logger.AddHook(hook)

				logger.Error("error")
			},
			func(fields logrus.Fields) {
				require.Equal(t, "reportcaller_test.go#58", fields["Location"])
			},
		)
	})

	t.Run("Location_SetKey", func(t *testing.T) {
		logrusutil.TestJSON(
			t,
			func(logger *logrus.Logger) {
				rootPath, err := os.Getwd()
				require.Nil(t, err)

				hook := NewReportCallerHook(rootPath)

				hook.SetKey("File")

				logger.AddHook(hook)

				logger.Error("error")
			},
			func(fields logrus.Fields) {
				require.Equal(t, "reportcaller_test.go:79", fields["File"])
			},
		)
	})
}
