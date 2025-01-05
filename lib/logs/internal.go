package logs

import (
	"fmt"
	"io"
	"time"
)

type internalLogger struct {
	writer io.Writer
	ignoreDebugLogs bool
}

var _ Logger = internalLogger{}

func (receiver internalLogger) log(level string, tag string, msg ...any) {
	if receiver.ignoreDebugLogs && level == levelDebug {
		return
	}

	var now string = time.Now().Format(time.RFC3339)

	var text string = fmt.Sprint(msg...)

	fmt.Fprintf(receiver.Writer(), "[%s] %s [%s] %s\n", tag, now, level, text)
}

func (receiver internalLogger) Writer() io.Writer {
	var writer io.Writer = receiver.writer
	if nil == writer {
		writer = io.Discard
	}
	return writer
}

func (receiver internalLogger) Debug(tag string, msg ...any) {
	const level = levelDebug
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Error(tag string, msg ...any) {
	const level = levelError
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Info(tag string, msg ...any) {
	const level = levelInfo
	receiver.log(level, tag, msg...)
}
