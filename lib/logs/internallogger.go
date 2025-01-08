package logs

type internalLogger struct {
	tag string
	internalTaggedLogger internalTaggedLogger
}

var _ Logger = internalLogger{}

func (receiver internalLogger) Debug(msg ...any) {
	receiver.internalTaggedLogger.Debug(receiver.tag, msg...)
}

func (receiver internalLogger) Debugf(format string, msg ...any) {
	receiver.internalTaggedLogger.Debugf(receiver.tag, format, msg...)
}

func (receiver internalLogger) Error(msg ...any) {
	receiver.internalTaggedLogger.Error(receiver.tag, msg...)
}

func (receiver internalLogger) Errorf(format string, msg ...any) {
	receiver.internalTaggedLogger.Errorf(receiver.tag, format, msg...)
}

func (receiver internalLogger) Info(msg ...any) {
	receiver.internalTaggedLogger.Info(receiver.tag, msg...)
}

func (receiver internalLogger) Infof(format string, msg ...any) {
	receiver.internalTaggedLogger.Infof(receiver.tag, format, msg...)
}
