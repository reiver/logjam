package logsrv

import (
	"github.com/mmcomp/go-log"

	"os"
)

var (
	Logger log.Logger = log.Logger{
		Output: os.Stdout,	
	}
)

func Begin(a ...interface{}) log.Logger {
	return Logger.Begin(a...)
}
