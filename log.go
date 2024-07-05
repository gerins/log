package log

import (
	"context"
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
	LevelFatal   = slog.Level(12)
	LevelTrace   = slog.Level(16) // Trace service to service communication
	LevelRequest = slog.Level(17) // Request log
)

var (
	globalLogger            *slog.Logger
	enableHideSensitiveData bool

	DefaultConfig = Config{
		LogToTerminal:     true,
		LogToFile:         false,
		Location:          "/log/",
		FileLogName:       "server_log",
		FileFormat:        ".%Y-%b-%d-%H-%M.log",
		MaxAge:            30,
		RotationFile:      24,
		Level:             LevelDebug,
		CustomWriter:      nil,
		HideSensitiveData: false,
	}
)

type (
	Config struct {
		LogToTerminal     bool       // Set log output to stdout
		LogToFile         bool       // Set log output to file
		Location          string     // Location file log will be save. Default "project_directory/log/".
		FileLogName       string     // File log name. Default "server_log".
		FileFormat        string     // Default "FileLogName.2021-Oct-22-00-00.log"
		MaxAge            int        // Days before deleting log file. Default 30 days.
		RotationFile      int        // Hour before creating new file. Default 24 hour.
		Level             slog.Level // Log output level. Default level DEBUG
		CustomWriter      io.Writer  // Specify custom writer for log output
		HideSensitiveData bool       // Enable hide sensitive data with struct tag `log:"hide"`
	}
)

func Init() {
	InitWithConfig(DefaultConfig)
}

func InitWithConfig(cfg Config) {
	if cfg.Location == "" {
		cfg.Location = DefaultConfig.Location
	}
	if cfg.FileLogName == "" {
		cfg.FileLogName = DefaultConfig.FileLogName
	}
	if cfg.FileFormat == "" {
		cfg.FileFormat = DefaultConfig.FileFormat
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = DefaultConfig.MaxAge
	}
	if cfg.RotationFile == 0 {
		cfg.RotationFile = DefaultConfig.RotationFile
	}
	if cfg.Level == 0 {
		cfg.Level = DefaultConfig.Level
	}

	enableHideSensitiveData = cfg.HideSensitiveData

	var output []io.Writer

	if cfg.LogToTerminal {
		output = append(output, os.Stdout)
	}

	if cfg.LogToFile {
		// Find current project directory
		currentDirectory, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed find current Directory for setup Log file, %s", err.Error())
		}

		fileLocation := currentDirectory + cfg.Location
		fileFormat := fmt.Sprintf("%s%s", cfg.FileLogName, cfg.FileFormat)

		// Initiate file rotate log
		fileWriter, err := rotatelogs.New(
			fileLocation+fileFormat,
			rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge)*24*time.Hour),                // Maximum time before deleting file log
			rotatelogs.WithRotationTime(time.Duration(cfg.RotationFile)*time.Hour),       // Time before creating new file
			rotatelogs.WithClock(rotatelogs.Local),                                       // Use local time for file rotation
			rotatelogs.WithLinkName(fmt.Sprintf("%s.log", fileLocation+cfg.FileLogName))) // Use file shortcut for accessing log file
		if err != nil {
			log.Fatalf("failed initiate file log rotation, %s", err.Error())
		}

		output = append(output, fileWriter)
	}

	if cfg.CustomWriter != nil {
		output = append(output, cfg.CustomWriter)
	}

	globalLogger = slog.New(slog.NewJSONHandler(io.MultiWriter(output...), &slog.HandlerOptions{
		Level: cfg.Level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove field msg if the value is empty
			if a.Key == slog.MessageKey && a.Value.String() == "" {
				return slog.Attr{}
			}

			// Customize the name of the level key and the output string, including custom level values.
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
				case level < LevelFatal:
					a.Value = slog.StringValue("ERROR")
				case level < LevelTrace:
					a.Value = slog.StringValue("FATAL")
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
	globalLogger.Log(context.Background(), LevelFatal, formatMultipleArguments(i))
	os.Exit(1)
}

func Fatalf(msg string, i ...any) {
	globalLogger.Log(context.Background(), LevelFatal, fmt.Sprintf(msg, i...))
	os.Exit(1)
}
