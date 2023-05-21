package log

import (
	"fmt"
	"math/rand"
	"runtime"

	"go.uber.org/zap/zapcore"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// formatMultipleArguments is used for formatting multiple argument input
func formatMultipleArguments(args []interface{}) string {
	var format string
	for i := range args {
		if i > 0 {
			format += " " // Add space between argument
		}
		format += "%+v" // When printing structs, the plus flag (%+v) adds field names
	}
	return fmt.Sprintf(format, args...)
}

// GetCaller return trimmed location of the file who call the log function
// example project_name/usecase/user.go:34
func GetCaller(level string, skip int) string {
	entryCaller := zapcore.NewEntryCaller(runtime.Caller(skip))
	return fmt.Sprintf("[%s] %s", level, entryCaller.TrimmedPath())
}
