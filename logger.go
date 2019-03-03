package zaplogger

import "go.uber.org/zap"

const (
	logKeyLoggerName = "logger.name"
)

func loggerName(name string) zap.Field {
	return zap.String(logKeyLoggerName, name)
}

// NewLogger creates a new zap.Logger based on root logger.
func NewLogger(root *zap.Logger, name string) *zap.Logger {
	return root.With(loggerName(name))
}
