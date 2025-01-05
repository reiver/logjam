package logsrv

import (
	"github.com/reiver/logjam/cfg"
	"github.com/reiver/logjam/lib/logs"
)

var Logger logs.TaggedLogger = logs.NewStdOutLogger(cfg.Config.ProdMode())

func Debug(tag string, msg ...any) {
	Logger.Debug(tag, msg...)
}

func Debugf(tag string, format string, msg ...any) {
	Logger.Debugf(tag, format, msg...)
}

func Error(tag string, msg ...any) {
	Logger.Error(tag, msg...)
}

func Errorf(tag string, format string, msg ...any) {
	Logger.Errorf(tag, format, msg...)
}

func Info(tag string, msg ...any) {
	Logger.Info(tag, msg...)
}

func Infof(tag string, format string, msg ...any) {
	Logger.Infof(tag, format, msg...)
}

func Tag(tag string) logs.Logger {
	return Logger.Tag(tag)
}
