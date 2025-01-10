package logs

import (
	"fmt"
	"io"
	"time"
)

type internalTaggedLogger struct {
	writer io.Writer
	ignoreDebugLogs bool
}

var _ TaggedLogger = internalTaggedLogger{}

func (receiver internalTaggedLogger) discard(level string) bool {
	return receiver.ignoreDebugLogs && level == levelDebug
}

func (receiver internalTaggedLogger) log(level string, tag string, msg ...any) {
	var text string = fmt.Sprint(msg...)

	receiver.publish(level, tag, text)
}

func (receiver internalTaggedLogger) logf(level string, tag string, format string, msg ...any) {
	var text string = fmt.Sprintf(format, msg...)

	receiver.publish(level, tag, text)
}

// publish is what actually writes the log
func (receiver internalTaggedLogger) publish(level string, tag string, text string) {
	if receiver.discard(level) {
		return
	}

	var now string = time.Now().Format(time.RFC3339)

	fmt.Fprintf(receiver.Writer(), "[%s] %s [%s] %s\n", tag, now, level, text)
}

func (receiver internalTaggedLogger) Writer() io.Writer {
	var writer io.Writer = receiver.writer
	if nil == writer {
		writer = io.Discard
	}
	return writer
}

func (receiver internalTaggedLogger) Debug(tag string, msg ...any) {
	const level = levelDebug
	receiver.log(level, tag, msg...)
}

func (receiver internalTaggedLogger) Debugf(tag string, format string, msg ...any) {
	const level = levelDebug
	receiver.logf(level, tag, format, msg...)
}

func (receiver internalTaggedLogger) Error(tag string, msg ...any) {
	const level = levelError
	receiver.log(level, tag, msg...)
}

func (receiver internalTaggedLogger) Errorf(tag string, format string, msg ...any) {
	const level = levelError
	receiver.logf(level, tag, format, msg...)
}

func (receiver internalTaggedLogger) Info(tag string, msg ...any) {
	const level = levelInfo
	receiver.log(level, tag, msg...)
}

func (receiver internalTaggedLogger) Infof(tag string, format string, msg ...any) {
	const level = levelInfo
	receiver.logf(level, tag, format, msg...)
}

func (receiver internalTaggedLogger) Tag(tag string) Logger {
	return internalLogger{
		tag:tag,
		internalTaggedLogger:receiver,
	}
}
