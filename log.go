package log

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	processID = "ProcessID"
)

var (
	globalLogger  *zap.Logger
	tracingLogger *zap.Logger

	DefaultLogConfig = Config{
		LogToTerminal: true,
		Location:      "/log/",
		FileLogName:   "server_log",
		FileFormat:    ".%Y-%b-%d-%H-%M.log",
		MaxAge:        30,
		RotationFile:  24,
		UseStackTrace: false,
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
		UseStackTrace bool   // Print stack trace to log file. Default false.
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

	// Initiate file rotate log
	rotateLog, err := rotatelogs.New(
		fileLocation+fileFormat,
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour),          // Maximum time before deleting file log
		rotatelogs.WithRotationTime(time.Duration(cfg.RotationFile)*time.Hour), // Time before creating new file
		rotatelogs.WithClock(rotatelogs.Local),                                 // Use local time for file rotation
		rotatelogs.WithLinkName(fileLocation+cfg.FileLogName))                  // Use file shortcut for accessing log file
	if err != nil {
		log.Fatalf("failed initiate file log rotation, %s", err.Error())
	}

	// Creating log encoder
	encoder := zap.NewProductionEncoderConfig()
	encoder.TimeKey = "time"
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder

	// Initiate log output
	outputToFile := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.AddSync(rotateLog), zapcore.InfoLevel)
	outputToTeriminal := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)

	// Building the zap core
	core := zapcore.NewTee(outputToFile)
	if cfg.LogToTerminal {
		core = zapcore.NewTee(outputToFile, outputToTeriminal)
	}

	// Creating Zap core
	if cfg.UseStackTrace {
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(1))
		tracingLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel), zap.AddCallerSkip(3))
	} else {
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		tracingLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(3))
	}
}

func Debug(i ...interface{}) {
	globalLogger.Debug(formatMultipleArguments(i))
}

func Debugf(format string, i ...interface{}) {
	globalLogger.Debug(fmt.Sprintf(format, i...))
}

func Info(i ...interface{}) {
	globalLogger.Info(formatMultipleArguments(i))
}

func Infof(format string, i ...interface{}) {
	globalLogger.Info(fmt.Sprintf(format, i...))
}

func Warn(i ...interface{}) {
	globalLogger.Warn(formatMultipleArguments(i))
}

func Warnf(format string, i ...interface{}) {
	globalLogger.Warn(fmt.Sprintf(format, i...))
}

func Error(i ...interface{}) {
	globalLogger.Error(formatMultipleArguments(i))
}

func Errorf(format string, i ...interface{}) {
	globalLogger.Error(fmt.Sprintf(format, i...))
}

func Fatal(i ...interface{}) {
	globalLogger.Fatal(formatMultipleArguments(i))
}

func Fatalf(msg string, i ...interface{}) {
	globalLogger.Fatal(fmt.Sprintf(msg, i...))
}

// Special log only for tracing service to service communication
func Tracing(processId, url, method string, resCode int, resPayload []byte, reqHeader, payload, respHeader interface{}, err error, dur int64) {
	var responsePayload interface{}
	if err := json.Unmarshal(resPayload, &responsePayload); err != nil {
		responsePayload = string(resPayload)
	}

	tracingLogger.Info("TRACING_LOG",
		zap.String(processID, processId),
		zap.String("URL", url),
		zap.String("Method", method),
		zap.Any("RequestHeader", reqHeader),
		zap.Any("RequestBody", payload),
		zap.Any("ResponseHeader", respHeader),
		zap.Any("ResponseBody", responsePayload),
		zap.Int("StatusCode", resCode),
		zap.Int64("RequestDuration", dur),
		zap.Any("Error", err),
	)
}
