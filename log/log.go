package log

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelCritical
)

// SetGCPLogger is middleware to set logger with trace in the context.
// In this middleware, projectID is required.
// If you don't set projectID, this middleware panics.
func SetGCPLogger(projectID, label string) func(*gin.Context) {
	return func(c *gin.Context) {
		if projectID == "" {
			panic("invalid projectID")
			return
		}
		ctx, done := SetSpanContext(c.Request, label)
		defer done()

		zctx := zerolog.New(os.Stdout).With()

		sc := GetSpanContext(ctx)
		if sc.TraceID != "" {
			trace := fmt.Sprintf("projects/%s/traces/%s", projectID, sc.TraceID)
			zctx = zctx.Str("logging.googleapis.com/trace", trace)
		}
		if sc.SpanID != "" {
			zctx = zctx.Str("logging.googleapis.com/spanId", sc.SpanID)
		}

		logger := zctx.Str("service", label).Logger()
		ctx = logger.WithContext(ctx)
		*c.Request = *c.Request.WithContext(ctx)
	}
}

func Debugf(ctx context.Context, format string, a ...interface{}) {
	loggingf(ctx, LogLevelDebug, format, a...)
}

func Infof(ctx context.Context, format string, a ...interface{}) {
	loggingf(ctx, LogLevelInfo, format, a...)
}

func Warningf(ctx context.Context, format string, a ...interface{}) {
	loggingf(ctx, LogLevelWarn, format, a...)
}

func Errorf(ctx context.Context, format string, a ...interface{}) {
	loggingf(ctx, LogLevelError, format, a...)
}

func Criticalf(ctx context.Context, format string, a ...interface{}) {
	loggingf(ctx, LogLevelCritical, format, a...)
}

func DebugObj(ctx context.Context, label string, obj interface{}) {
	loggingObj(ctx, LogLevelDebug, label, obj)
}

func InfoObj(ctx context.Context, label string, obj interface{}) {
	loggingObj(ctx, LogLevelInfo, label, obj)
}

func WarningObj(ctx context.Context, label string, obj interface{}) {
	loggingObj(ctx, LogLevelWarn, label, obj)
}

func ErrorObj(ctx context.Context, label string, obj interface{}) {
	loggingObj(ctx, LogLevelError, label, obj)
}

func CriticalObj(ctx context.Context, label string, obj interface{}) {
	loggingObj(ctx, LogLevelCritical, label, obj)
}

func loggingf(ctx context.Context, level LogLevel, format string, a ...interface{}) {
	if level >= LogLevelError {
		zlog.Ctx(ctx).
			Log().
			Str("severity", toSeverity(level)).
			Str("stack_trace", fmt.Sprintf(format, a...)+":\n"+callers().String()). // messageを含まないと同じ関数の別のエラーが同一と見なされてしまうため
			Msgf(format, a...)
	} else {
		zlog.Ctx(ctx).
			Log().
			Str("severity", toSeverity(level)).
			Msgf(format, a...)
	}
}

func loggingObj(ctx context.Context, level LogLevel, label string, obj interface{}) {
	if level >= LogLevelError {
		zlog.Ctx(ctx).
			Log().
			Str("severity", toSeverity(level)).
			Str("stack_trace", fmt.Sprintf("%#v", obj)+":\n"+callers().String()). // messageを含まないと同じ関数の別のエラーが同一と見なされてしまうため
			Interface(label, obj).
			Msgf("%#v", obj)
	} else {
		zlog.Ctx(ctx).
			Log().
			Str("severity", toSeverity(level)).
			Interface(label, obj).
			Msgf("%#v", obj)
	}
}

func toSeverity(level LogLevel) string {
	var severity string
	switch level {
	case LogLevelDebug:
		severity = "DEBUG"
	case LogLevelWarn:
		severity = "WARNING"
	case LogLevelError:
		severity = "ERROR"
	case LogLevelCritical:
		severity = "CRITICAL"
	default:
		severity = "INFO"
	}
	return severity
}
