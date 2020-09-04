package log

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
)

const contextLoggerKey = "logger"

// DefaultFormatter sets up the default key names for a logrus Logger
var DefaultFormatter = &logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyTime:  "timestamp",
		logrus.FieldKeyMsg:   "message",
		logrus.FieldKeyLevel: "level",
	},
}

// Opt is an option that can be passed to New to configure the logger to be used in the app.
type Opt func(*logrus.Logger)

// WithFormatter returns an Opt that sets the formatter on a logger
func WithFormatter(formatter logrus.Formatter) Opt {
	return func(logger *logrus.Logger) {
		logger.Formatter = formatter
	}
}

// WithLevel returns an Opt that sets the log level on a logger
func WithLevel(level logrus.Level) Opt {
	return func(logger *logrus.Logger) {
		logger.Level = level
	}
}

// WithOutput returns an Opt that sets the output for a logger
func WithOutput(out io.Writer) Opt {
	return func(logger *logrus.Logger) {
		logger.Out = out
	}
}

// ContextWithLogger injects a logger into the context for future reuse
func ContextWithLogger(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger)
}

// FromContext fetches a logger from the context. This function is not lenient, it gets really angry if the context
// has no logger :)
func FromContext(ctx context.Context) logrus.FieldLogger {
	l, ok := ctx.Value(contextLoggerKey).(logrus.FieldLogger)
	if !ok {
		panic("unset logger")
	}

	return l
}

// New creates a logger with the desired options
func New(opts ...Opt) *logrus.Logger {
	logger := logrus.New()
	for _, opt := range opts {
		opt(logger)
	}

	return logger
}
