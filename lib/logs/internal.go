package logs

import (
	"fmt"
	"io"
	"strings"
	"time"
)

type internalLogger struct {
	writer io.Writer
	ignoreDebugLogs bool
}

var _ Logger = internalLogger{}

func (receiver internalLogger) log(level string, tag string, msg ...string) {
	if receiver.ignoreDebugLogs && level == levelDebug {
		return
	}
	fmt.Fprintf(receiver.Writer(), "[%s] %s [%s] %s\n", tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", "))
}

func (receiver internalLogger) Writer() io.Writer {
	var writer io.Writer = receiver.writer
	if nil == writer {
		writer = io.Discard
	}
	return writer
}

func (receiver internalLogger) Debug(tag string, msg ...string) {
	const level = levelDebug
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Error(tag string, msg ...string) {
	const level = levelError
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Info(tag string, msg ...string) {
	const level = levelInfo
	receiver.log(level, tag, msg...)
}
