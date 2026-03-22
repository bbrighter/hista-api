package hista

import (
	"context"
	"time"

	"encore.dev/rlog"
	"gorm.io/gorm/logger"
)

type EncoreLogger struct {
	level logger.LogLevel
}

func (l EncoreLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.level = level
	return l
}

func (l EncoreLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Info {
		rlog.Info(msg, args...)
	}
}

func (l EncoreLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Warn {
		rlog.Warn(msg, args...)
	}
}

func (l EncoreLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.level >= logger.Error {
		rlog.Error(msg, args...)
	}
}

func (l EncoreLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.level <= logger.Silent {
		return
	}
	sql, rows := fc()
	duration := time.Since(begin)
	if err != nil {
		rlog.Error("SQL error", "sql", sql, "rows", rows, "duration", duration, "error", err)
	} else {
		rlog.Info("SQL executed", "sql", sql, "rows", rows, "duration", duration)
	}
}
