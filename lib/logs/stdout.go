package logs

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type stdoutLogger struct {
	ignoreDebugLogs bool
}

var _ Logger = stdoutLogger{}

func (receiver stdoutLogger) log(tag string, level string, msg ...string) error {
	if receiver.ignoreDebugLogs && level == levelDebug {
		return nil
	}
	fmt.Fprintf(os.Stdout, "[%s] %s [%s] %s\n", tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", "))
	return nil
}

func (receiver stdoutLogger) Debug(tag string, msg ...string) error {
	const level = levelDebug
	return receiver.log(tag, level, msg...)
}

func (receiver stdoutLogger) Error(tag string, msg ...string) error {
	const level = levelError
	return receiver.log(tag, level, msg...)
}

func (receiver stdoutLogger) Info(tag string, msg ...string) error {
	const level = levelInfo
	return receiver.log(tag, level, msg...)
}

func NewStdOutLogger(ignoreDebugLogs bool) Logger {
	return stdoutLogger{
		ignoreDebugLogs: ignoreDebugLogs,
	}
}
