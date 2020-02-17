package logger

import (
	"github.com/sirupsen/logrus"
)

// Level log level
type Level = logrus.Level

// Fields arbitrary fields
type Fields = logrus.Fields

// FieldLogger ...
type FieldLogger = logrus.FieldLogger

// Entry ...
type Entry = logrus.Entry

// Logger ...
type Logger struct {
	entry *Entry
}

// NewLogger ...
func NewLogger() *Logger {
	l := logrus.New()

	l.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	return &Logger{entry: logrus.NewEntry(l)}
}

// Debug logs a message at level Debug on the standard logger.
func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

// Debugln logs a message at level Debug on the standard logger.
func (l *Logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

// Infoln logs a message at level Info on the standard logger.
func (l *Logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

// Infof logs a message at level Info on the standard logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (l *Logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

// Warnln logs a message at level Warn on the standard logger.
func (l *Logger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

// Warnf logs a message at level Warn on the standard logger.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

// Errorln logs a message at level Error on the standard logger.
func (l *Logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger.
func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

// Fatalln logs a message at level Fatal on the standard logger.
func (l *Logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

// Fatalf logs a message at level Fatal on the standard logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

// SetLevel sets the level of the logger
func (l *Logger) SetLevel(level Level) {
	l.entry.Logger.SetLevel(level)
}

// WithFields add fields to next log
func (l *Logger) WithFields(fields Fields) *Entry {
	return l.entry.WithFields(fields)
}
