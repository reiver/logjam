package logs

type Logger interface {
	Debug(msg ...any)
	Debugf(format string, msg ...any)
	Error(msg ...any)
	Errorf(format string, msg ...any)
	Info(msg ...any)
	Infof(format string, msg ...any)
}
