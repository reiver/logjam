package logs

type Logger interface {
	Debug(tag string, msg ...string)
	Error(tag string, msg ...string)
	Info(tag string, msg ...string)
}
