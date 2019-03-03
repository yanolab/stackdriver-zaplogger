package zaplogger

import (
	"encoding/json"

	"cloud.google.com/go/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewCore creates a new Core instance.
func NewCore(cli *logging.Client, lv zapcore.LevelEnabler) *Core {
	return &Core{
		LevelEnabler: lv,
		cli:          cli,
		logger:       cli.Logger("default"),
		enc:          zapcore.NewJSONEncoder(EncoderConfig),
	}
}

// Core is a stack driver core.
type Core struct {
	zapcore.LevelEnabler
	cli    *logging.Client
	logger *logging.Logger
	enc    zapcore.Encoder
}

// With recreates a new Core with given fields.
func (c *Core) With(fields []zapcore.Field) zapcore.Core {
	core := c.clone()

	for _, f := range fields {
		if f.Key == logKeyLoggerName {
			core.logger = c.cli.Logger(f.String)
			continue
		}
		f.AddTo(core.enc)
	}

	return core
}

// Check checks log level.
func (c *Core) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return ce.AddCore(entry, c)
	}

	return ce
}

// Write writes zap entry and fields to stackdrvier.
func (c *Core) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	e, err := c.toStackdriver(entry, fields)
	if err != nil {
		return err
	}
	c.logger.Log(e)
	return nil
}

// Sync syncs stack driver logger.
func (c *Core) Sync() error {
	return c.logger.Flush()
}

func (c *Core) clone() *Core {
	return &Core{
		cli:          c.cli,
		logger:       c.logger,
		enc:          c.enc.Clone(),
		LevelEnabler: c.LevelEnabler,
	}
}

func (c *Core) toStackdriver(entry zapcore.Entry, fields []zapcore.Field) (logging.Entry, error) {
	serverity, ok := logLevelSeverity[entry.Level]
	if !ok {
		serverity = logging.Default
	}

	fields = append(fields, zap.String("stack", entry.Stack))
	b, err := c.enc.EncodeEntry(entry, fields)
	if err != nil {
		return logging.Entry{}, err
	}
	var payload interface{}
	if err := json.Unmarshal(b.Bytes(), &payload); err != nil {
		return logging.Entry{}, err
	}

	return logging.Entry{
		Timestamp: entry.Time,
		Severity:  serverity,
		LogName:   entry.LoggerName,
		Payload:   payload,
	}, nil
}
