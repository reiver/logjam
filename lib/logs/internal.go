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

func (receiver internalLogger) discard(level string) bool {
	return receiver.ignoreDebugLogs && level == levelDebug
}

func (receiver internalLogger) log(level string, tag string, msg ...any) {
	var text string = fmt.Sprint(msg...)

	receiver.publish(level, tag, text)
}

func (receiver internalLogger) logf(level string, tag string, format string, msg ...any) {
	var text string = fmt.Sprintf(format, msg...)

	receiver.publish(level, tag, text)
}

// publish is what actually writes the log
func (receiver internalLogger) publish(level string, tag string, text string) {
	if receiver.discard(level) {
		return
	}

	var now string = time.Now().Format(time.RFC3339)

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

func (receiver internalLogger) Debugf(tag string, format string, msg ...any) {
	const level = levelDebug
	receiver.logf(level, tag, format, msg...)
}

func (receiver internalLogger) Error(tag string, msg ...any) {
	const level = levelError
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Errorf(tag string, format string, msg ...any) {
	const level = levelError
	receiver.logf(level, tag, format, msg...)
}

func (receiver internalLogger) Info(tag string, msg ...any) {
	const level = levelInfo
	receiver.log(level, tag, msg...)
}

func (receiver internalLogger) Infof(tag string, format string, msg ...any) {
	const level = levelInfo
	receiver.logf(level, tag, format, msg...)
}
