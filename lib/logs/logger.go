package logs

type Logger interface {
	Debug(tag string, msg ...any)
	Debugf(tag string, format string, msg ...any)
	Error(tag string, msg ...any)
	Errorf(tag string, format string, msg ...any)
	Info(tag string, msg ...any)
	Infof(tag string, format string, msg ...any)
}
