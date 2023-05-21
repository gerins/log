package gorm

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"

	"github.com/gerins/log"
)

const (
	callerName      = "DATABASE"
	callerSkipLevel = 3
)

var ErrRecordNotFound = errors.New("record not found")

type Config struct {
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
	LogLevel                  logger.LogLevel
}

var (
	Default = New(Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Info,
	})
)

func New(config Config) logger.Interface {
	var (
		// Log format
		infoStr      = "[INFO] %s "
		warnStr      = "[WARN] %s "
		errStr       = "[ERROR] %s "
		traceStr     = "[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s [%.3fms] [rows:%v] %s"
	)

	return &logExtension{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logExtension struct {
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logExtension) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l logExtension) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		log.Context(ctx).SubLog(l.getCallerLocation(), fmt.Sprintf(l.infoStr+msg, data...))
	}
}

// Warn print warn messages
func (l logExtension) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		log.Context(ctx).SubLog(l.getCallerLocation(), fmt.Sprintf(l.warnStr+msg, data...))
	}
}

// Error print error messages
func (l logExtension) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		log.Context(ctx).SubLog(l.getCallerLocation(), fmt.Sprintf(l.errStr+msg, data...))
	}
}

// Trace print sql message
func (l logExtension) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	level := l.getCallerLocation()
	duration := float64(elapsed.Nanoseconds()) / 1e6

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceErrStr, err, duration, "-", sql))
		} else {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceErrStr, err, duration, rows, sql))
		}

	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceWarnStr, slowLog, duration, "-", sql))
		} else {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceWarnStr, slowLog, duration, rows, sql))
		}

	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceStr, duration, "-", sql))
		} else {
			log.Context(ctx).SubLog(level, fmt.Sprintf(l.traceStr, duration, rows, sql))
		}
	}
}

func (l logExtension) getCallerLocation() string {
	for i := callerSkipLevel; i < 15; i++ {
		filePath := zapcore.NewEntryCaller(runtime.Caller(i)).TrimmedPath()
		if !strings.HasPrefix(filePath, "gorm@") {
			return fmt.Sprintf("[%s] %s", callerName, filePath)
		}
	}
	return ""
}
