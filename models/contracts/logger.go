package contracts

type TLogLevel string

const (
	LInfo  TLogLevel = "info"
	LDebug TLogLevel = "debug"
	LError TLogLevel = "error"
)

type ILogger interface {
	Log(tag string, level TLogLevel, msg ...string) error
}
