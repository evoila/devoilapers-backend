package util

import (
	"OperatorAutomation/pkg/utils/logger"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ReflectionLogger_Error(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RError(errors.New("Some err"), "Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 1, len(hookEntries))

	// Check values
	assert.Equal(t, "Some msg", hookEntries[0].Message)

	errMsg := fmt.Sprintf("%s", hookEntries[0].Data["err"])
	assert.Equal(t, "Some err", errMsg)

	lineInfo := hookEntries[0].Data["pos"]
	assert.Equal(t, "reflection_log_test.go:OperatorAutomation/test/unit_tests/util.Test_ReflectionLogger_Error:18", lineInfo)
}

func Test_ReflectionLogger_Warning(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.WarnLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RWarn("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 1, len(hookEntries))

	// Check values
	assert.Equal(t, "Some msg", hookEntries[0].Message)

	lineInfo := hookEntries[0].Data["pos"]
	assert.Equal(t, "reflection_log_test.go:OperatorAutomation/test/unit_tests/util.Test_ReflectionLogger_Warning:37", lineInfo)
}

func Test_ReflectionLogger_Info(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RInfo("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 1, len(hookEntries))

	// Check values
	assert.Equal(t, "Some msg", hookEntries[0].Message)

	lineInfo := hookEntries[0].Data["pos"]
	assert.Equal(t, "reflection_log_test.go:OperatorAutomation/test/unit_tests/util.Test_ReflectionLogger_Info:53", lineInfo)
}

func Test_ReflectionLogger_Trace(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.TraceLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RTrace("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 1, len(hookEntries))

	// Check values
	assert.Equal(t, "Some msg", hookEntries[0].Message)

	lineInfo := hookEntries[0].Data["pos"]
	assert.Equal(t, "reflection_log_test.go:OperatorAutomation/test/unit_tests/util.Test_ReflectionLogger_Trace:69", lineInfo)
}

func Test_ReflectionLogger_Warning_WrongLevel(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RWarn("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 0, len(hookEntries))
}

func Test_ReflectionLogger_Info_WrongLevel(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.WarnLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RInfo("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 0, len(hookEntries))
}

func Test_ReflectionLogger_Trace_WrongLevel(t *testing.T) {
	_, hook := test.NewNullLogger()
	logrus.SetLevel(logrus.InfoLevel)
	logrus.AddHook(hook)
	// Produce error
	logger.RTrace("Some msg")
	hookEntries := hook.Entries
	assert.Equal(t, 0, len(hookEntries))
}
