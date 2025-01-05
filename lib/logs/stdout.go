package logs

import (
	"os"
)

// NewStdOutLogger returns a logger that write logs to Stadard-Out (stdout).
//
// See also [NewLogger]
func NewStdOutLogger(ignoreDebugLogs bool) TaggedLogger {
	return NewLogger(os.Stdout, ignoreDebugLogs)
}
