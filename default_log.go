package colorlog

import "io"

var defaultLog = New()

func Fatal(text string) {
	defaultLog.Fatal(text)
}

func Fatalf(format string, args ...any) {
	defaultLog.Fatalf(format, args...)
}

func Error(text string) {
	defaultLog.Error(text)
}

func Errorf(format string, args ...any) {
	defaultLog.Errorf(format, args...)
}

func Warn(text string) {
	defaultLog.Warn(text)
}

func Warnf(format string, args ...any) {
	defaultLog.Warnf(format, args...)
}

func Info(text string) {
	defaultLog.Info(text)
}

func Infof(format string, args ...any) {
	defaultLog.Infof(format, args...)
}

func Debug(text string) {
	defaultLog.Debug(text)
}

func Debugf(format string, args ...any) {
	defaultLog.Debugf(format, args...)
}

func WithConfig(cfg *Config) {
	defaultLog.WithConfig(cfg)
}

func WithOut(iow io.Writer) {
	defaultLog.WithOut(iow)
}

func WithMask(mask string) {
	defaultLog.WithMask(mask)
}
