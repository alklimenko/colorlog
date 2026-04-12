package colorlog

import "io"

var defaultLog = New()

func Error(text string, args ...any) {
	defaultLog.Error(text, args...)
}

func Warn(text string, args ...any) {
	defaultLog.Warn(text, args...)
}

func Info(text string, args ...any) {
	defaultLog.Info(text, args...)
}

func Debug(text string, args ...any) {
	defaultLog.Debug(text, args...)
}

func WithConfig(cfg *Config) {
	defaultLog.WithConfig(cfg)
}

func WithErr(iow io.Writer) {
	defaultLog.WithErr(iow)
}

func WithOut(iow io.Writer) {
	defaultLog.WithOut(iow)
}

func WithMask(mask string) {
	defaultLog.WithMask(mask)
}
