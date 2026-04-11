package colorlog_test

import (
	"testing"

	"github.com/alklimenko/colorlog"
)

func TestLog(t *testing.T) {
	var l = colorlog.New().WithDefaultLightConfig()

	l.Fatal("Fatal message")
	l.Error("Error message")
	l.Warn("Warning message")
	l.Info("Info message")
	l.Debug("Debug message")
}

func TestLog2(t *testing.T) {
	colorlog.Fatal("Fatal message")
	colorlog.Error("Error message")
	colorlog.Warn("Warning message")
	colorlog.Info("Info message")
	colorlog.Debug("Debug message")
}
