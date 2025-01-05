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

func (receiver stdoutLogger) log(tag string, level string, msg ...string) {
	if receiver.ignoreDebugLogs && level == levelDebug {
		return
	}
	fmt.Fprintf(os.Stdout, "[%s] %s [%s] %s\n", tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", "))
}

func (receiver stdoutLogger) Debug(tag string, msg ...string) {
	const level = levelDebug
	receiver.log(tag, level, msg...)
}

func (receiver stdoutLogger) Error(tag string, msg ...string) {
	const level = levelError
	receiver.log(tag, level, msg...)
}

func (receiver stdoutLogger) Info(tag string, msg ...string) {
	const level = levelInfo
	receiver.log(tag, level, msg...)
}

func NewStdOutLogger(ignoreDebugLogs bool) Logger {
	return stdoutLogger{
		ignoreDebugLogs: ignoreDebugLogs,
	}
}
