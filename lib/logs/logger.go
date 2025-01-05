package logs

type Logger interface {
	Debug(tag string, msg ...string) error
	Error(tag string, msg ...string) error
	Info(tag string, msg ...string) error
}
