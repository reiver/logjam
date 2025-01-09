package db

import (
	"github.com/reiver/logjam/lib/logs"
)

type Context struct {
	Logger logs.Logger
}

var _ logs.Logger = Context{}

func (receiver Context) Debug(msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Debug(msg...)
}

func (receiver Context) Debugf(format string, msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Debugf(format, msg...)
}

func (receiver Context) Error(msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Error(msg...)
}

func (receiver Context) Errorf(format string, msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Errorf(format, msg...)
}

func (receiver Context) Info(msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Info(msg...)
}

func (receiver Context) Infof(format string, msg ...any) {
	if nil == receiver.Logger {
		return
	}

	receiver.Logger.Infof(format, msg...)
}
