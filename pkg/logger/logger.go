package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	zap *zap.SugaredLogger
}

var defaultLogger *Logger

func init() {
	defaultLogger = New()
}

func New() *Logger {
	config := zap.NewDevelopmentConfig()

	zapLogger, err := config.Build()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	return &Logger{
		zap: zapLogger.Sugar(),
	}
}

func NewProduction() *Logger {
	config := zap.NewProductionConfig()

	zapLogger, err := config.Build()
	if err != nil {
		panic("Failed to initialize production logger: " + err.Error())
	}

	return &Logger{
		zap: zapLogger.Sugar(),
	}
}

func (l *Logger) Info(v ...any) {
	l.zap.Info(v...)
}

func (l *Logger) Infof(format string, v ...any) {
	l.zap.Infof(format, v...)
}

func (l *Logger) Error(v ...any) {
	l.zap.Error(v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.zap.Errorf(format, v...)
}

func (l *Logger) Debug(v ...any) {
	l.zap.Debug(v...)
}

func (l *Logger) Debugf(format string, v ...any) {
	l.zap.Debugf(format, v...)
}

func (l *Logger) Warn(v ...any) {
	l.zap.Warn(v...)
}

func (l *Logger) Warnf(format string, v ...any) {
	l.zap.Warnf(format, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.zap.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.zap.Fatalf(format, v...)
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() {
	l.zap.Sync()
}

// Global functions for convenience (keeping the same wrapper APIs)
func Info(v ...any) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}

func Error(v ...any) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}

func Debug(v ...any) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

func Warn(v ...any) {
	defaultLogger.Warn(v...)
}

func Warnf(format string, v ...any) {
	defaultLogger.Warnf(format, v...)
}

func Fatal(v ...any) {
	defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...any) {
	defaultLogger.Fatalf(format, v...)
}

// Sync flushes any buffered log entries - useful for graceful shutdown
func Sync() {
	defaultLogger.Sync()
}
