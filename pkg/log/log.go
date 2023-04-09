package log

import (
	"chino/models"
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field struct {
	Key   string
	Value any
}

func (f *Field) ToZap() zapcore.Field {
	return zap.Any(f.Key, f.Value)
}

type Logger struct {
	l *zap.Logger
}

func New() (*Logger, error) {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "msg",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		return nil, err
	}
	return &Logger{l: logger}, nil
}

func (l *Logger) Raw() *zap.Logger {
	return l.l
}

func (l *Logger) prepareFields(ctx context.Context, fields ...Field) []zap.Field {
	f := []zapcore.Field{}
	if pk := ctx.Value("user"); pk != nil {
		f = append(f, zap.Any("user", pk.(string)))
	}

	if reqID := ctx.Value("requestid"); reqID != nil {
		f = append(f, zap.Any("request_id", reqID.(string)))
	}
	for _, fd := range fields {
		f = append(f, fd.ToZap())
	}
	return f
}

func get(ctx context.Context) *Logger {
	return ctx.Value(models.String("logger")).(*Logger)
}

func (l *Logger) debug(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Debug(msg, f...)
}

func (l *Logger) info(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Info(msg, f...)
}

func (l *Logger) warn(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Warn(msg, f...)
}

func (l *Logger) error(ctx context.Context, msg string, fields ...Field) {
	f := l.prepareFields(ctx, fields...)
	l.l.Error(msg, f...)
}

func Debug(ctx context.Context, msg string, fields ...Field) {
	get(ctx).debug(ctx, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...Field) {
	get(ctx).info(ctx, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...Field) {
	get(ctx).warn(ctx, msg, fields...)
}

func Error(ctx context.Context, err error, fields ...Field) {
	get(ctx).error(ctx, err.Error(), fields...)
}
