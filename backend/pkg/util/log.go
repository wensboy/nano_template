package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogFormat controls how entries are serialized before being written.
type LogFormat string

const (
	LogFormatConsole  LogFormat = "console"
	LogFormatJSON     LogFormat = "json"
	LogFormatTemplate LogFormat = "template"
)

const (
	defaultLogTemplate = "{{time}} {{level|color}} {{caller|underline}} {{message}} {{fields|json}}"
	defaultTimeFormat  = "2006-01-02T15:04:05.000Z07:00"
)

// LoggerConfig describes an independent logger instance.
//
// Template syntax uses {{token|modifier}} expressions. Built-in tokens are:
// time, level, logger, caller, message, stacktrace, fields and field:<name>.
// Useful modifiers include upper, lower, json, kv, color, bold and underline.
type LoggerConfig struct {
	Level           zapcore.Level
	Format          LogFormat
	Template        string
	Writers         []io.Writer
	TimeFormat      string
	UseColor        bool
	DisableCaller   bool
	CallerSkip      int
	StacktraceLevel zapcore.LevelEnabler
}

// Logger is a small zap-backed logger with pluggable output and formatting.
type Logger struct {
	zap *zap.Logger
}

var (
	globalLoggerMu sync.RWMutex
	globalLogger   *Logger
)

// NewLogger creates an independent logger. When no writer is supplied it logs to stdout.
func NewLogger(cfg LoggerConfig) (*Logger, error) {
	cfg = normalizeLoggerConfig(cfg)

	ws := buildWriteSyncer(cfg.Writers)
	level := zap.NewAtomicLevelAt(cfg.Level)
	core := buildCore(cfg, ws, level)

	options := []zap.Option{zap.AddCallerSkip(1 + cfg.CallerSkip)}
	if !cfg.DisableCaller {
		options = append(options, zap.AddCaller())
	}
	if cfg.StacktraceLevel != nil {
		options = append(options, zap.AddStacktrace(cfg.StacktraceLevel))
	}

	return &Logger{zap: zap.New(core, options...)}, nil
}

// MustNewLogger creates a logger and panics on configuration errors.
func MustNewLogger(cfg LoggerConfig) *Logger {
	logger, err := NewLogger(cfg)
	if err != nil {
		panic(err)
	}
	return logger
}

// InitLogger initializes the package-level logger used by Info, Warn, Error and Debug.
func InitLogger(logToFile bool, filePath string) {
	writers := []io.Writer{os.Stdout}
	if logToFile {
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		writers = append(writers, file)
	}

	SetDefaultLogger(MustNewLogger(LoggerConfig{
		Level:           zapcore.DebugLevel,
		Format:          LogFormatTemplate,
		Writers:         writers,
		UseColor:        true,
		CallerSkip:      1,
		StacktraceLevel: zapcore.ErrorLevel,
	}))
}

// SetDefaultLogger replaces the package-level logger used by the helper functions.
func SetDefaultLogger(logger *Logger) {
	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()
	globalLogger = logger
}

// Zap exposes the underlying zap logger for integrations that need zap-specific APIs.
func (l *Logger) Zap() *zap.Logger {
	if l == nil || l.zap == nil {
		return zap.NewNop()
	}
	return l.zap
}

// Named returns a child logger with a zap logger name.
func (l *Logger) Named(name string) *Logger {
	return &Logger{zap: l.Zap().Named(name)}
}

// With returns a child logger carrying the supplied structured fields.
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{zap: l.Zap().With(fields...)}
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Zap().Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Zap().Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Zap().Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Zap().Error(msg, fields...)
}

func (l *Logger) Sync() error {
	return l.Zap().Sync()
}

// Info logs an info message through the package-level logger.
func Info(msg string, fields ...zap.Field) {
	defaultLogger().Info(msg, fields...)
}

// Error logs an error message through the package-level logger.
func Error(msg string, fields ...zap.Field) {
	defaultLogger().Error(msg, fields...)
}

// Debug logs a debug message through the package-level logger.
func Debug(msg string, fields ...zap.Field) {
	defaultLogger().Debug(msg, fields...)
}

// Warn logs a warning message through the package-level logger.
func Warn(msg string, fields ...zap.Field) {
	defaultLogger().Warn(msg, fields...)
}

// Sync flushes buffered log entries from the package-level logger.
func Sync() {
	_ = defaultLogger().Sync()
}

func defaultLogger() *Logger {
	globalLoggerMu.RLock()
	logger := globalLogger
	globalLoggerMu.RUnlock()
	if logger != nil {
		return logger
	}

	globalLoggerMu.Lock()
	defer globalLoggerMu.Unlock()
	if globalLogger == nil {
		globalLogger = MustNewLogger(LoggerConfig{
			Level:           zapcore.DebugLevel,
			Format:          LogFormatTemplate,
			Writers:         []io.Writer{os.Stdout},
			UseColor:        true,
			CallerSkip:      1,
			StacktraceLevel: zapcore.ErrorLevel,
		})
	}
	return globalLogger
}

func normalizeLoggerConfig(cfg LoggerConfig) LoggerConfig {
	if cfg.Format == "" {
		cfg.Format = LogFormatTemplate
	}
	if cfg.Template == "" {
		cfg.Template = defaultLogTemplate
	}
	if cfg.TimeFormat == "" {
		cfg.TimeFormat = defaultTimeFormat
	}
	if len(cfg.Writers) == 0 {
		cfg.Writers = []io.Writer{os.Stdout}
	}
	return cfg
}

func buildCore(cfg LoggerConfig, ws zapcore.WriteSyncer, level zapcore.LevelEnabler) zapcore.Core {
	if cfg.Format == LogFormatTemplate {
		return &templateCore{
			level:    level,
			ws:       ws,
			template: cfg.Template,
			timeFmt:  cfg.TimeFormat,
			useColor: cfg.UseColor,
		}
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if cfg.Format == LogFormatJSON {
		return zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), ws, level)
	}
	return zapcore.NewCore(zapcore.NewConsoleEncoder(encoderCfg), ws, level)
}

func buildWriteSyncer(writers []io.Writer) zapcore.WriteSyncer {
	syncers := make([]zapcore.WriteSyncer, 0, len(writers))
	for _, writer := range writers {
		if writer == nil {
			continue
		}
		syncer, ok := writer.(zapcore.WriteSyncer)
		if !ok {
			syncer = zapcore.AddSync(writer)
		}
		syncers = append(syncers, zapcore.Lock(syncer))
	}
	if len(syncers) == 0 {
		return zapcore.AddSync(io.Discard)
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

type templateCore struct {
	level    zapcore.LevelEnabler
	ws       zapcore.WriteSyncer
	template string
	timeFmt  string
	useColor bool
	fields   []zap.Field
}

func (c *templateCore) Enabled(level zapcore.Level) bool {
	return c.level.Enabled(level)
}

func (c *templateCore) With(fields []zap.Field) zapcore.Core {
	clone := *c
	clone.fields = append(clone.fields[:len(clone.fields):len(clone.fields)], fields...)
	return &clone
}

func (c *templateCore) Check(entry zapcore.Entry, checked *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(entry.Level) {
		return checked.AddCore(entry, c)
	}
	return checked
}

func (c *templateCore) Write(entry zapcore.Entry, fields []zap.Field) error {
	rendered, err := c.render(entry, fields)
	if err != nil {
		return err
	}
	_, err = c.ws.Write([]byte(rendered))
	return err
}

func (c *templateCore) Sync() error {
	return c.ws.Sync()
}

func (c *templateCore) render(entry zapcore.Entry, fields []zap.Field) (string, error) {
	allFields := make([]zap.Field, 0, len(c.fields)+len(fields))
	allFields = append(allFields, c.fields...)
	allFields = append(allFields, fields...)

	encodedFields := encodeFields(allFields)
	ctx := renderContext{
		entry:   entry,
		fields:  encodedFields,
		timeFmt: c.timeFmt,
		color:   c.useColor,
	}

	var out strings.Builder
	template := c.template
	for {
		start := strings.Index(template, "{{")
		if start < 0 {
			out.WriteString(template)
			break
		}
		out.WriteString(template[:start])
		template = template[start+2:]

		end := strings.Index(template, "}}")
		if end < 0 {
			out.WriteString("{{")
			out.WriteString(template)
			break
		}

		token := strings.TrimSpace(template[:end])
		out.WriteString(ctx.renderToken(token))
		template = template[end+2:]
	}

	if entry.Stack != "" {
		rendered := out.String()
		if !strings.Contains(rendered, entry.Stack) {
			out.WriteByte('\n')
			out.WriteString(entry.Stack)
		}
	}
	out.WriteByte('\n')
	return out.String(), nil
}

type renderContext struct {
	entry   zapcore.Entry
	fields  map[string]interface{}
	timeFmt string
	color   bool
}

func (c renderContext) renderToken(expr string) string {
	parts := strings.Split(expr, "|")
	if len(parts) == 0 {
		return ""
	}

	token := strings.TrimSpace(parts[0])
	value := c.resolve(token)
	for _, rawModifier := range parts[1:] {
		modifier := strings.TrimSpace(rawModifier)
		if token == "fields" {
			switch modifier {
			case "json":
				value = renderFieldsJSON(c.fields)
				continue
			case "kv":
				value = renderFieldsKV(c.fields)
				continue
			}
		}
		if strings.HasPrefix(token, "field:") && modifier == "json" {
			name := strings.TrimSpace(strings.TrimPrefix(token, "field:"))
			value = renderValueJSON(c.fields[name])
			continue
		}
		value = c.applyModifier(value, modifier)
	}
	return value
}

func (c renderContext) resolve(token string) string {
	switch {
	case token == "time":
		return c.entry.Time.Format(c.timeFmt)
	case token == "level":
		return c.entry.Level.String()
	case token == "logger":
		return c.entry.LoggerName
	case token == "caller":
		if !c.entry.Caller.Defined {
			return ""
		}
		return c.entry.Caller.TrimmedPath()
	case token == "message" || token == "msg":
		return c.entry.Message
	case token == "stacktrace" || token == "stack":
		return c.entry.Stack
	case token == "fields":
		return renderFieldsJSON(c.fields)
	case strings.HasPrefix(token, "field:"):
		name := strings.TrimSpace(strings.TrimPrefix(token, "field:"))
		return renderSingleField(c.fields[name])
	default:
		return ""
	}
}

func (c renderContext) applyModifier(value string, modifier string) string {
	switch modifier {
	case "":
		return value
	case "upper":
		return strings.ToUpper(value)
	case "lower":
		return strings.ToLower(value)
	case "json":
		return renderValueJSON(value)
	case "kv":
		return renderFieldsKV(c.fields)
	case "color":
		if !c.color {
			return value
		}
		return colorizeLevel(c.entry.Level, value)
	case "bold":
		if !c.color {
			return value
		}
		return "\x1b[1m" + value + "\x1b[0m"
	case "underline":
		if !c.color {
			return value
		}
		return "\x1b[4m" + value + "\x1b[0m"
	default:
		return value
	}
}

func encodeFields(fields []zap.Field) map[string]interface{} {
	encoder := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(encoder)
	}
	return encoder.Fields
}

func renderFieldsJSON(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return "{}"
	}
	data, err := json.Marshal(fields)
	if err != nil {
		return fmt.Sprint(fields)
	}
	return string(data)
}

func renderValueJSON(value interface{}) string {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprint(value)
	}
	return string(data)
}

func renderFieldsKV(fields map[string]interface{}) string {
	if len(fields) == 0 {
		return ""
	}

	keys := make([]string, 0, len(fields))
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	values := make([]string, 0, len(keys))
	for _, key := range keys {
		values = append(values, fmt.Sprintf("%s=%s", key, renderSingleField(fields[key])))
	}
	return strings.Join(values, " ")
}

func renderSingleField(value interface{}) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case time.Time:
		return typed.Format(time.RFC3339Nano)
	case fmt.Stringer:
		return typed.String()
	case error:
		return typed.Error()
	default:
		switch typed := typed.(type) {
		case []byte:
			return string(typed)
		default:
			data, err := json.Marshal(typed)
			if err == nil {
				return string(data)
			}
			return fmt.Sprint(typed)
		}
	}
}

func colorizeLevel(level zapcore.Level, value string) string {
	color := "\x1b[37m"
	switch level {
	case zapcore.DebugLevel:
		color = "\x1b[36m"
	case zapcore.InfoLevel:
		color = "\x1b[32m"
	case zapcore.WarnLevel:
		color = "\x1b[33m"
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = "\x1b[31m"
	}
	return color + value + "\x1b[0m"
}

var _ zapcore.Core = (*templateCore)(nil)
