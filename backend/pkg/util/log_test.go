package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestTemplateLoggerRendersStructuredFields(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{level|upper}}|{{message}}|{{fields|json}}",
		Writers:  []io.Writer{&buf},
	})

	logger.Info("created", zap.String("user", "ada"), zap.Int("count", 2))

	want := "INFO|created|{\"count\":2,\"user\":\"ada\"}\n"
	if got := buf.String(); got != want {
		t.Fatalf("unexpected template output:\nwant: %q\n got: %q", want, got)
	}
}

func TestTemplateLoggerSupportsKVAndSingleFieldTokens(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{message}}|{{field:request_id}}|{{field:count|json}}|{{fields|kv}}",
		Writers:  []io.Writer{&buf},
	})

	logger.With(zap.String("app", "nano")).Info(
		"accepted",
		zap.String("request_id", "req-1"),
		zap.Int("count", 3),
	)

	want := "accepted|req-1|3|app=nano count=3 request_id=req-1\n"
	if got := buf.String(); got != want {
		t.Fatalf("unexpected template output:\nwant: %q\n got: %q", want, got)
	}
}

func TestLoggersAreIndependentInstances(t *testing.T) {
	var first bytes.Buffer
	var second bytes.Buffer

	firstLogger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{message}}",
		Writers:  []io.Writer{&first},
	})
	secondLogger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{message}}",
		Writers:  []io.Writer{&second},
	})

	firstLogger.Info("first")
	secondLogger.Info("second")

	if got := first.String(); got != "first\n" {
		t.Fatalf("first logger wrote %q", got)
	}
	if got := second.String(); got != "second\n" {
		t.Fatalf("second logger wrote %q", got)
	}
}

func TestTemplateLoggerSupportsTerminalStyles(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{level|color}} {{message|bold}} {{caller|underline}}",
		Writers:  []io.Writer{&buf},
		UseColor: true,
	})

	logger.Info("styled")

	got := buf.String()
	for _, want := range []string{
		"\x1b[32minfo\x1b[0m",
		"\x1b[1mstyled\x1b[0m",
		"\x1b[4m",
		"log_test.go:",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("styled output %q does not contain %q", got, want)
		}
	}
}

func TestJSONLoggerKeepsZapStructuredOutput(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:   zapcore.DebugLevel,
		Format:  LogFormatJSON,
		Writers: []io.Writer{&buf},
	})

	logger.Warn("json message", zap.String("trace_id", "req-1"))

	var payload map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &payload); err != nil {
		t.Fatalf("failed to unmarshal JSON log: %v\nlog: %s", err, buf.String())
	}

	assertJSONField(t, payload, "level", "warn")
	assertJSONField(t, payload, "msg", "json message")
	assertJSONField(t, payload, "trace_id", "req-1")

	caller, ok := payload["caller"].(string)
	if !ok || !strings.Contains(caller, "log_test.go:") {
		t.Fatalf("caller was not recorded from the test callsite: %#v", payload["caller"])
	}
}

func TestCallerPointsAtDirectLogCall(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.DebugLevel,
		Format:   LogFormatTemplate,
		Template: "{{caller}}",
		Writers:  []io.Writer{&buf},
	})

	_, _, line, _ := runtime.Caller(0)
	logger.Info("caller")

	want := fmt.Sprintf("log_test.go:%d", line+1)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("caller mismatch: want %q in %q", want, got)
	}
	if strings.Contains(buf.String(), "log.go:") {
		t.Fatalf("caller should not point at the logger wrapper: %q", buf.String())
	}
}

func TestCallerSkipSkipsAdditionalWrapper(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:      zapcore.DebugLevel,
		Format:     LogFormatTemplate,
		Template:   "{{caller}}",
		Writers:    []io.Writer{&buf},
		CallerSkip: 1,
	})

	wrapper := func() {
		logger.Info("from wrapper")
	}

	_, _, line, _ := runtime.Caller(0)
	wrapper()

	want := fmt.Sprintf("log_test.go:%d", line+1)
	if got := buf.String(); !strings.Contains(got, want) {
		t.Fatalf("caller skip mismatch: want %q in %q", want, got)
	}
}

func TestPackageLevelLoggerCompatibility(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:      zapcore.DebugLevel,
		Format:     LogFormatTemplate,
		Template:   "{{caller}} {{message}} {{fields|json}}",
		Writers:    []io.Writer{&buf},
		CallerSkip: 1,
	})

	globalLoggerMu.Lock()
	old := globalLogger
	globalLogger = logger
	globalLoggerMu.Unlock()
	t.Cleanup(func() {
		globalLoggerMu.Lock()
		globalLogger = old
		globalLoggerMu.Unlock()
	})

	_, _, line, _ := runtime.Caller(0)
	Info("global", zap.Bool("ok", true))

	wantCaller := fmt.Sprintf("log_test.go:%d", line+1)
	got := buf.String()
	if !strings.Contains(got, wantCaller) {
		t.Fatalf("global helper caller mismatch: want %q in %q", wantCaller, got)
	}
	if !strings.Contains(got, "global {\"ok\":true}") {
		t.Fatalf("global helper did not keep structured fields: %q", got)
	}
}

func TestLevelFiltering(t *testing.T) {
	var buf bytes.Buffer
	logger := mustTestLogger(t, LoggerConfig{
		Level:    zapcore.WarnLevel,
		Format:   LogFormatTemplate,
		Template: "{{level}} {{message}}",
		Writers:  []io.Writer{&buf},
	})

	logger.Info("hidden")
	logger.Warn("visible")

	if got := buf.String(); got != "warn visible\n" {
		t.Fatalf("unexpected level filtering result: %q", got)
	}
}

func mustTestLogger(t *testing.T, cfg LoggerConfig) *Logger {
	t.Helper()

	logger, err := NewLogger(cfg)
	if err != nil {
		t.Fatalf("NewLogger failed: %v", err)
	}
	return logger
}

func assertJSONField(t *testing.T, payload map[string]interface{}, key string, want string) {
	t.Helper()

	got, ok := payload[key].(string)
	if !ok {
		t.Fatalf("JSON field %q is not a string: %#v", key, payload[key])
	}
	if got != want {
		t.Fatalf("JSON field %q mismatch: want %q, got %q", key, want, got)
	}
}
