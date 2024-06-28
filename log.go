package log

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

const (
	processID    = "processID"
	levelDebug   = "DEBUG"
	levelInfo    = "INFO"
	levelWarn    = "WARN"
	levelError   = "ERROR"
	levelFatal   = "FATAL"
	levelRequest = "REQUEST_LOG"
	levelTracing = "TRACING_LOG"
)

var (
	globalLogger *slog.Logger

	DefaultLogConfig = Config{
		LogToTerminal: true,
		Location:      "/log/",
		FileLogName:   "server_log",
		FileFormat:    ".%Y-%b-%d-%H-%M.log",
		MaxAge:        30,
		RotationFile:  24,
	}
)

type (
	Config struct {
		LogToTerminal bool   // Highly recommended disable in production server. Default true.
		Location      string // Location file log will be save. Default "project_directory/log/".
		FileLogName   string // File log name. Default "server_log".
		FileFormat    string // Default "FileLogName.2021-Oct-22-00-00.log"
		MaxAge        int    // Days before deleting log file. Default 30 days.
		RotationFile  int    // Hour before creating new file. Default 24 hour.
	}
)

func Init() {
	InitWithConfig(DefaultLogConfig)
}

func InitWithConfig(cfg Config) {
	if cfg.Location == "" {
		cfg.Location = DefaultLogConfig.Location
	}
	if cfg.FileLogName == "" {
		cfg.FileLogName = DefaultLogConfig.FileLogName
	}
	if cfg.FileFormat == "" {
		cfg.FileFormat = DefaultLogConfig.FileFormat
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = DefaultLogConfig.MaxAge
	}
	if cfg.RotationFile == 0 {
		cfg.RotationFile = DefaultLogConfig.RotationFile
	}

	// Find current project directory
	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed find current Directory for setup Log file, %s", err.Error())
	}

	fileLocation := currentDirectory + cfg.Location
	fileFormat := fmt.Sprintf("%s%s", cfg.FileLogName, cfg.FileFormat)

	var output io.Writer

	// Initiate file rotate log
	output, err = rotatelogs.New(
		fileLocation+fileFormat,
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour),                // Maximum time before deleting file log
		rotatelogs.WithRotationTime(time.Duration(cfg.RotationFile)*time.Hour),       // Time before creating new file
		rotatelogs.WithClock(rotatelogs.Local),                                       // Use local time for file rotation
		rotatelogs.WithLinkName(fmt.Sprintf("%s.log", fileLocation+cfg.FileLogName))) // Use file shortcut for accessing log file
	if err != nil {
		log.Fatalf("failed initiate file log rotation, %s", err.Error())
	}

	if cfg.LogToTerminal {
		output = io.MultiWriter(os.Stdout, output)
	}

	globalLogger = slog.New(slog.NewJSONHandler(output, nil))
}

func Debug(i ...any) {
	globalLogger.Debug(formatMultipleArguments(i))
}

func Debugf(format string, i ...any) {
	globalLogger.Debug(fmt.Sprintf(format, i...))
}

func Info(i ...any) {
	globalLogger.Info(formatMultipleArguments(i))
}

func Infof(format string, i ...any) {
	globalLogger.Info(fmt.Sprintf(format, i...))
}

func Warn(i ...any) {
	globalLogger.Warn(formatMultipleArguments(i))
}

func Warnf(format string, i ...any) {
	globalLogger.Warn(fmt.Sprintf(format, i...))
}

func Error(i ...any) {
	globalLogger.Error(formatMultipleArguments(i))
}

func Errorf(format string, i ...any) {
	globalLogger.Error(fmt.Sprintf(format, i...))
}

func Fatal(i ...any) {
	globalLogger.Error(formatMultipleArguments(i))
	os.Exit(1)
}

func Fatalf(msg string, i ...any) {
	globalLogger.Error(fmt.Sprintf(msg, i...))
	os.Exit(1)
}

// Trace service to service communication
func Tracing(processId, url, method string, resCode int, resPayload []byte, reqHeader, payload, respHeader any, err error, dur int64) {
	var responsePayload any
	if err := json.Unmarshal(resPayload, &responsePayload); err != nil {
		responsePayload = string(resPayload)
	}
	go globalLogger.LogAttrs(context.Background(), slog.LevelInfo, levelTracing,
		slog.String("caller", GetCaller("", 3)),
		slog.String(processID, processId),
		slog.String("method", method),
		slog.String("url", url),
		slog.Int("statusCode", resCode),
		slog.Int64("requestDuration", dur),
		slog.Any("error", err),
		slog.Any("requestHeader", reqHeader),
		slog.Any("requestBody", payload),
		slog.Any("responseHeader", respHeader),
		slog.Any("responseBody", responsePayload),
	)
}
