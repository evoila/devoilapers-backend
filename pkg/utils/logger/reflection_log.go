package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

// Error log with line information
func RError(err error, msg ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		logrus.WithFields(
			logrus.Fields{
				"err": err,
				"pos": fmt.Sprintf("%s:%s:%d", file, funcName, line),
			}).Error(msg...)
	} else {
		logrus.Error(msg)
	}
}

// Warning log with line information
func RWarn(msg ...interface{}) {
	if !logrus.IsLevelEnabled(logrus.WarnLevel) {
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		logrus.WithFields(
			logrus.Fields{
				"pos": fmt.Sprintf("%s:%s:%d", file, funcName, line),
			}).Warn(msg...)
	}else {
		logrus.Warn(msg)
	}
}

// Info log with line information
func RInfo(msg ...interface{}) {
	if !logrus.IsLevelEnabled(logrus.InfoLevel) {
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		logrus.WithFields(
			logrus.Fields{
				"pos": fmt.Sprintf("%s:%s:%d", file, funcName, line),
			}).Info(msg...)
	} else {
		logrus.Info(msg)
	}
}

// Trace log with line information
func RTrace(msg ...interface{}) {
	if !logrus.IsLevelEnabled(logrus.TraceLevel) {
		return
	}

	if pc, file, line, ok := runtime.Caller(1); ok {
		file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
		logrus.WithFields(
			logrus.Fields{
				"pos": fmt.Sprintf("%s:%s:%d", file, funcName, line),
			}).Trace(msg...)
	} else {
		logrus.Trace(msg)
	}
}
