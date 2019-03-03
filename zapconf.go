package zaplogger

import (
	"cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

var logLevelSeverity = map[zapcore.Level]logging.Severity{
	zapcore.DebugLevel:  logging.Debug,
	zapcore.InfoLevel:   logging.Info,
	zapcore.WarnLevel:   logging.Warning,
	zapcore.ErrorLevel:  logging.Error,
	zapcore.DPanicLevel: logging.Critical,
	zapcore.PanicLevel:  logging.Alert,
	zapcore.FatalLevel:  logging.Emergency,
}

// EncoderConfig is default encode config for stack driver.
var EncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "eventTime",
	LevelKey:       "severity",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    encodeLevel,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func encodeLevel(lv zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	severity := logLevelSeverity[lv]
	enc.AppendString(severity.String())
}
