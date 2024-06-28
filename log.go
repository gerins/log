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
	processID = "processID"

	// Logging level from least important to most important
	LevelDebug   = slog.LevelDebug
	LevelInfo    = slog.LevelInfo
	LevelWarning = slog.LevelWarn
	LevelError   = slog.LevelError
	LevelTrace   = slog.Level(16) // Trace service to service communication
	LevelRequest = slog.Level(17) // Request log
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
		Level:         LevelDebug,
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
		Level         slog.Level
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
	if cfg.Level == 0 {
		cfg.Level = DefaultLogConfig.Level
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

	globalLogger = slog.New(slog.NewJSONHandler(output, &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize the name of the level key and the output string, including
			// custom level values.
			if a.Key == slog.LevelKey {
				// Handle custom level values.
				level := a.Value.Any().(slog.Level)
				switch {
				case level < LevelInfo:
					a.Value = slog.StringValue("DEBUG")
				case level < LevelWarning:
					a.Value = slog.StringValue("INFO")
				case level < LevelError:
					a.Value = slog.StringValue("WARN")
				case level < LevelTrace:
					a.Value = slog.StringValue("ERROR")
				case level < LevelRequest:
					a.Value = slog.StringValue("TRACE")
				default:
					a.Value = slog.StringValue("REQUEST")
				}
			}
			return a
		},
	}))
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
	go globalLogger.LogAttrs(context.Background(), LevelTrace, "",
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
