package logs

type Logger interface {
	Debug(tag string, msg ...any)
	Error(tag string, msg ...any)
	Info(tag string, msg ...any)
}
