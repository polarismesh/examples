package logx

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/tal-tech/go-zero/core/timex"
	"go.opentelemetry.io/otel/trace"
)

type traceLogger struct {
	logEntry
	Trace string `json:"trace,omitempty"`
	Span  string `json:"span,omitempty"`
	ctx   context.Context
}

func (l *traceLogger) Error(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprint(v...), durationCallerDepth))
	}
}

func (l *traceLogger) Errorf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, formatWithCaller(fmt.Sprintf(format, v...), durationCallerDepth))
	}
}

func (l *traceLogger) Errorv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(errorLog, levelError, v)
	}
}

func (l *traceLogger) Info(v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Infof(format string, v ...interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Infov(v interface{}) {
	if shallLog(InfoLevel) {
		l.write(infoLog, levelInfo, v)
	}
}

func (l *traceLogger) Slow(v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprint(v...))
	}
}

func (l *traceLogger) Slowf(format string, v ...interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, fmt.Sprintf(format, v...))
	}
}

func (l *traceLogger) Slowv(v interface{}) {
	if shallLog(ErrorLevel) {
		l.write(slowLog, levelSlow, v)
	}
}

func (l *traceLogger) WithDuration(duration time.Duration) Logger {
	l.Duration = timex.ReprOfDuration(duration)
	return l
}

func (l *traceLogger) write(writer io.Writer, level string, val interface{}) {
	l.Timestamp = getTimestamp()
	l.Level = level
	l.Content = val
	l.Trace = traceIdFromContext(l.ctx)
	l.Span = spanIdFromContext(l.ctx)
	outputJson(writer, l)
}

// WithContext sets ctx to log, for keeping tracing information.
func WithContext(ctx context.Context) Logger {
	return &traceLogger{
		ctx: ctx,
	}
}

func spanIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}

	return ""
}

func traceIdFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}

	return ""
}
