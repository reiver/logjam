package logsrv

import (
	"github.com/sparkscience/logjam/backend/arg"

	"github.com/mmcomp/go-log"

	"os"
)

var (
	Logger log.Logger = log.Logger{
		Output: os.Stdout,
	}
)

func init() {
	switch {
	case arg.VeryVeryVeryVeryVeryVerbose:
		Logger = Logger.Level(6)
	case arg.VeryVeryVeryVeryVerbose:
		Logger = Logger.Level(5)
	case arg.VeryVeryVeryVerbose:
		Logger = Logger.Level(4)
	case arg.VeryVeryVerbose:
		Logger = Logger.Level(3)
	case arg.VeryVerbose:
		Logger = Logger.Level(2)
	case arg.Verbose:
		Logger = Logger.Level(1)
	default:
		Logger = Logger.Level(0)
	}
}

func Begin(a ...interface{}) log.Logger {
	return Logger.Begin(a...)
}
