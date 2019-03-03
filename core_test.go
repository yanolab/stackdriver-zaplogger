package zaplogger

import (
	"strings"
	"testing"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"
)

func Test_toStackdriver(t *testing.T) {
	core := &Core{
		LevelEnabler: zap.DebugLevel,
		enc:          zapcore.NewJSONEncoder(EncoderConfig),
	}

	zapentry := zapcore.Entry{Stack: "stack"}
	stdentry, err := core.toStackdriver(zapentry, []zapcore.Field{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if zapentry.Time.Unix() != stdentry.Timestamp.Unix() {
		t.Fatalf("unexpected time: expected=%v actual:%v", zapentry.Time, stdentry.Timestamp)
	}
	if zapentry.Level.CapitalString() != strings.ToUpper(stdentry.Severity.String()) {
		t.Fatalf("unexpected log level: expected=%v actual:%v", zapentry.Level, stdentry.Severity)
	}

	m := stdentry.Payload.(map[string]interface{})
	if v := m["stack"]; v != zapentry.Stack {
		t.Fatalf("unexpected stack: expected=%s actual:%s", zapentry.Stack, v)
	}
}
