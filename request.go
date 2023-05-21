package log

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

var (
	subLogSkipLevel = 2
	logRequestKey   contextKey
)

type (
	contextKey int

	// Data Model for tracking information of incoming request
	request struct {
		processID  string
		UserID     int
		IP         string
		Method     string
		URL        string
		ReqHeader  interface{}
		ReqBody    interface{}
		RespHeader interface{}
		RespBody   interface{}
		StatusCode int                    // HTTP status code
		timeStart  time.Time              // Capture when the request start
		ExtraData  map[string]interface{} // Additional data
		subLogs    []subLog               // Sub logging data
		WaitGroup  *sync.WaitGroup
	}

	// Data model for saving all log output in request flow
	subLog struct {
		Level   string `json:"level"`
		Message string `json:"message"`
	}
)

// NewRequest will create new log data model for incoming request
func NewRequest() *request {
	return &request{
		processID: generateRandomString(20),
		timeStart: time.Now(),
		ExtraData: make(map[string]interface{}),
		WaitGroup: new(sync.WaitGroup),
	}
}

// Save will save current request information to log file
func (m *request) Save() {
	go func() {
		m.WaitGroup.Wait() // Wait for all goroutine finish before logging
		globalLogger.Info("REQUEST_LOG",
			zap.String(processID, m.processID),
			zap.Int("UserID", m.UserID),
			zap.String("IP", m.IP),
			zap.String("Method", m.Method),
			zap.String("URL", m.URL),
			zap.Any("RequestHeader", m.ReqHeader),
			zap.Any("RequestBody", m.ReqBody),
			zap.Any("ResponseHeader", m.RespHeader),
			zap.Any("ResponseBody", m.RespBody),
			zap.Int("StatusCode", m.StatusCode),
			zap.Int64("RequestDuration", int64(time.Since(m.timeStart).Milliseconds())),
			zap.Any("ExtraData", m.ExtraData),
			zap.Any("SubLog", m.subLogs),
		)
	}()
}

// SetProcessID is used for set process id as your preferences format.
func (m *request) SetProcessID(value string) {
	m.processID = value
}

// ProcessID is used for get current process id from log request model.
func (m *request) ProcessID() string {
	return m.processID
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

func (m *request) Debug(i ...interface{}) {
	msg := formatMultipleArguments(i)
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("DEBUG", subLogSkipLevel), Message: msg})
}

func (m *request) Debugf(format string, i ...interface{}) {
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("DEBUG", subLogSkipLevel), Message: fmt.Sprintf(format, i...)})
}

func (m *request) Info(i ...interface{}) {
	msg := formatMultipleArguments(i)
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("INFO", subLogSkipLevel), Message: msg})
}

func (m *request) Infof(format string, i ...interface{}) {
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("INFO", subLogSkipLevel), Message: fmt.Sprintf(format, i...)})
}

func (m *request) Warn(i ...interface{}) {
	msg := formatMultipleArguments(i)
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("WARN", subLogSkipLevel), Message: msg})
}

func (m *request) Warnf(format string, i ...interface{}) {
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("WARN", subLogSkipLevel), Message: fmt.Sprintf(format, i...)})
}

func (m *request) Error(i ...interface{}) {
	msg := formatMultipleArguments(i)
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("ERROR", subLogSkipLevel), Message: msg})
}

func (m *request) Errorf(format string, i ...interface{}) {
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("ERROR", subLogSkipLevel), Message: fmt.Sprintf(format, i...)})
}

func (m *request) Fatal(i ...interface{}) {
	msg := formatMultipleArguments(i)
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("FATAL", subLogSkipLevel), Message: msg})
}

func (m *request) Fatalf(format string, i ...interface{}) {
	m.subLogs = append(m.subLogs, subLog{Level: GetCaller("FATAL", subLogSkipLevel), Message: fmt.Sprintf(format, i...)})
}

func (m *request) SubLog(level, message string) {
	m.subLogs = append(m.subLogs, subLog{Level: level, Message: message})
}
