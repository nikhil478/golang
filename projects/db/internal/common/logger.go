package common

import "go.uber.org/zap"

type logger struct {
	logger *zap.Logger
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
}

func NewLogger() Logger {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return &logger{
		logger: zapLogger,
	}
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *logger) With(fields ...zap.Field) Logger {
	return &logger{
		logger: l.logger.With(fields...),
	}
}
