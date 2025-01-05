package logger

import (
	"fmt"
	"github.com/reiver/logjam/models/contracts"
	"strings"
	"time"
)

type stdoutLogger struct {
	ignoreDebugLogs bool
}

func (s *stdoutLogger) Log(tag string, level contracts.TLogLevel, msg ...string) error {
	if s.ignoreDebugLogs && level == contracts.LDebug {
		return nil
	}
	println(fmt.Sprintf(`[%s] %s [%s] %s`, tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", ")))
	return nil
}

func NewSTDOUTLogger(ignoreDebugLogs bool) contracts.ILogger {
	return &stdoutLogger{
		ignoreDebugLogs: ignoreDebugLogs,
	}
}
