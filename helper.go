package log

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"strings"

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
func formatMultipleArguments(args []any) string {
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
	if level == "" {
		return entryCaller.TrimmedPath()
	}
	return fmt.Sprintf("[%s] %s", level, entryCaller.TrimmedPath())
}

func maskSensitiveData(payload any) {
	// Check if req is a pointer to a struct
	v := reflect.ValueOf(payload)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return
	}

	v = v.Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := t.Field(i).Tag.Get("log")
		if tag == "hide" && field.Kind() == reflect.String {
			field.SetString(strings.Repeat("*", len(field.String())))
		}
	}
}
