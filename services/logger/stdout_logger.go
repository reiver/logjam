package logger

import (
	"fmt"
	"github.com/sparkscience/logjam/models/contracts"
	"strings"
	"time"
)

type stdoutLogger struct {
}

func (s *stdoutLogger) Log(tag string, level contracts.TLogLevel, msg ...string) error {
	println(fmt.Sprintf(`[%s] %s [%s] %s`, tag, time.Now().Format(time.RFC3339), level, strings.Join(msg, ", ")))
	return nil
}

func NewSTDOUTLogger() contracts.ILogger {
	return &stdoutLogger{}
}
