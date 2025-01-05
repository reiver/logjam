package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/reiver/logjam/lib/logs"
)

type stdoutLogger struct {
	ignoreDebugLogs bool
}

func (s *stdoutLogger) Log(tag string, level logs.TLogLevel, msg ...string) error {
	if s.ignoreDebugLogs && level == logs.LDebug {
		return nil
	}
	println(fmt.Sprintf(`[%s] %s [%s] %s`, tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", ")))
	return nil
}

func NewSTDOUTLogger(ignoreDebugLogs bool) logs.ILogger {
	return &stdoutLogger{
		ignoreDebugLogs: ignoreDebugLogs,
	}
}
