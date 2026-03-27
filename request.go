package log

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
)

var (
	subLogSkipLevel = 2
	logRequestKey   contextKey
)

const (
	subLevelDebug = "DEBUG"
	subLevelInfo  = "INFO"
	subLevelWarn  = "WARN"
	subLevelError = "ERROR"
	subLevelFatal = "FATAL"
)

type (
	contextKey int

	// Data Model for tracking information of incoming request
	request struct {
		traceID    string
		IP         string
		Method     string
		URL        string
		ReqHeader  any
		ReqBody    any
		RespHeader any
		RespBody   any
		StatusCode int             // HTTP status code or other code
		timeStart  time.Time       // Capture when the request start
		ExtraData  map[string]any  // Additional data
		subLogs    []subLog        // Sub logging data
		WaitGroup  *sync.WaitGroup // Wait for all goroutine finish before printing log
	}

	// Data model for saving all log output in single request flow
	subLog struct {
		Level   string `json:"level"`
		Message string `json:"msg"`
	}
)

// NewRequest will create new log data model for incoming request
func NewRequest() *request {
	return &request{
		traceID:   generateRandomString(20),
		timeStart: time.Now(),
		ExtraData: make(map[string]any),
		WaitGroup: new(sync.WaitGroup),
	}
}

// Save will save current request information to log file
func (m *request) Save() {
	go func() {
		m.WaitGroup.Wait() // Wait for all goroutine finish before logging

		if enableHideSensitiveData {
			for _, data := range m.ExtraData {
				maskSensitiveData(data)
			}
			maskSensitiveData(m.ReqBody)
			maskSensitiveData(m.RespBody)
		}

		globalLogger.LogAttrs(context.Background(), LevelRequest, "",
			slog.String("caller", GetCaller("", 1)),
			slog.String(traceID, m.traceID),
			slog.String("ip", m.IP),
			slog.String("method", m.Method),
			slog.String("url", m.URL),
			slog.Int("statusCode", m.StatusCode),
			slog.Int64("totalDuration", time.Since(m.timeStart).Milliseconds()),
			slog.Any("requestHeader", m.ReqHeader),
			slog.Any("requestBody", m.ReqBody),
			slog.Any("responseHeader", m.RespHeader),
			slog.Any("responseBody", m.RespBody),
			slog.Any("extraData", m.ExtraData),
			slog.Any("subLog", m.subLogs),
		)
	}()
}

// SetTraceID is used for set trace id as your preferences format.
func (m *request) SetTraceID(traceID string) {
	m.traceID = traceID
}

// TraceID is used for get current process id from log request model.
func (m *request) TraceID() string {
	return m.traceID
}

func (m *request) SaveToContext(parent context.Context) context.Context {
	return context.WithValue(parent, logRequestKey, m)
}

// Context is used for get log request model from context
func Context(ctx context.Context) *request {
	data, ok := ctx.Value(logRequestKey).(*request)
	if !ok {
		data = NewRequest()
	}
	return data
}

// RecordDuration is used for record total duration a process could take
func (m *request) RecordDuration(processName string) processData {
	return processData{request: m, name: processName, timeStart: time.Now()}
}

func (m *request) Debug(i ...any) {
	msg := formatMultipleArguments(i)

	if disableSubLogs {
		m.globalLog(LevelDebug, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelDebug, subLogSkipLevel), Message: msg})
}

func (m *request) Debugf(format string, i ...any) {
	msg := fmt.Sprintf(format, i...)

	if disableSubLogs {
		m.globalLog(LevelDebug, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelDebug, subLogSkipLevel), Message: msg})
}

func (m *request) Info(i ...any) {
	msg := formatMultipleArguments(i)

	if disableSubLogs {
		m.globalLog(LevelInfo, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelInfo, subLogSkipLevel), Message: msg})
}

func (m *request) Infof(format string, i ...any) {
	msg := fmt.Sprintf(format, i...)

	if disableSubLogs {
		m.globalLog(LevelInfo, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelInfo, subLogSkipLevel), Message: msg})
}

func (m *request) Warn(i ...any) {
	msg := formatMultipleArguments(i)

	if disableSubLogs {
		m.globalLog(LevelWarning, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelWarn, subLogSkipLevel), Message: msg})
}

func (m *request) Warnf(format string, i ...any) {
	msg := fmt.Sprintf(format, i...)

	if disableSubLogs {
		m.globalLog(LevelWarning, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelWarn, subLogSkipLevel), Message: msg})
}

func (m *request) Error(i ...any) {
	msg := formatMultipleArguments(i)

	if disableSubLogs {
		m.globalLog(LevelError, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelError, subLogSkipLevel), Message: msg})
}

func (m *request) Errorf(format string, i ...any) {
	msg := fmt.Sprintf(format, i...)

	if disableSubLogs {
		m.globalLog(LevelError, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelError, subLogSkipLevel), Message: msg})
}

func (m *request) Fatal(i ...any) {
	msg := formatMultipleArguments(i)

	if disableSubLogs {
		m.globalLog(LevelFatal, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelFatal, subLogSkipLevel), Message: msg})
}

func (m *request) Fatalf(format string, i ...any) {
	msg := fmt.Sprintf(format, i...)

	if disableSubLogs {
		m.globalLog(LevelFatal, msg, "")
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: GetCaller(subLevelFatal, subLogSkipLevel), Message: msg})
}

func (m *request) SubLog(levelAndCaller, message string) {
	if disableSubLogs {
		m.globalLog(LevelInfo, message, levelAndCaller)
		return
	}

	m.subLogs = append(m.subLogs, subLog{Level: levelAndCaller, Message: message})
}

func (m *request) globalLog(level slog.Level, msg string, caller string) {
	if caller == "" {
		caller = GetCaller("", subLogSkipLevel+1)
	}

	attrs := []slog.Attr{
		slog.String("caller", caller),
		slog.String(traceID, m.traceID),
		slog.String("msg", msg),
	}

	globalLogger.LogAttrs(context.Background(), level, "", attrs...)
}
