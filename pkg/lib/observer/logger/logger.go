package logger

import "go.uber.org/zap/zapcore"

func (l *logger) Debug(log string, fields ...zapcore.Field) {
	l.log.Desugar().Debug(log, fields...)
}

func (l *logger) Info(log string, fields ...zapcore.Field) {
	l.log.Desugar().Info(log, fields...)
}

func (l *logger) Warning(log string, fields ...zapcore.Field) {
	l.log.Desugar().Warn(log, fields...)
}

func (l *logger) Error(log string, fields ...zapcore.Field) {
	l.log.Desugar().Error(log, fields...)
}
