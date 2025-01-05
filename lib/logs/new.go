package logs

import (
	"io"
)

// NewLogger returns a logger that write logs to the specified `writer`.
//
// See also [NewStdOutLogger]
func NewLogger(writer io.Writer, ignoreDebugLogs bool) TaggedLogger {
	return internalLogger{
		writer:writer,
		ignoreDebugLogs: ignoreDebugLogs,
	}
}
