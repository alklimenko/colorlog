package colorlog

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	TextRed          = "\033[31m"
	TextGreen        = "\033[32m"
	TextYellow       = "\033[33m"
	TextBlue         = "\033[34m"
	TextMagenta      = "\033[35m"
	TextCyan         = "\033[36m"
	TextGray         = "\033[37m"
	TextLightRed     = "\033[91m"
	TextLightGreen   = "\033[92m"
	TextLightYellow  = "\033[93m"
	TextLightBlue    = "\033[94m"
	TextLightMagenta = "\033[95m"
	TextLightCyan    = "\033[96m"
	TextWhite        = "\033[97m"
	TextReset        = "\033[0m"
	TextBold         = "\033[1m"
	TextItalic       = "\033[3m"
	TextUnder        = "\033[4m"
	TextBlink        = "\033[5m"
)

type ColorLog struct {
	errW io.Writer
	outW io.Writer
	cfg  *Config
	date time.Time
}

// Config - Ascii коды для установки соответствующих элементов лога
type Config struct {
	Date  string
	Time  string
	Fatal string
	Error string
	Warn  string
	Info  string
	Debug string
}

var (
	defaultDarkConfig = Config{
		Date:  TextBold + TextLightCyan,
		Time:  TextLightCyan,
		Fatal: TextBold + TextItalic + TextLightMagenta,
		Error: TextBold + TextLightRed,
		Warn:  TextBold + TextLightYellow,
		Info:  TextLightGreen,
		Debug: TextGray,
	}
	defaultLightConfig = Config{
		Date:  TextBold + TextCyan,
		Time:  TextCyan,
		Fatal: TextBold + TextItalic + TextLightMagenta,
		Error: TextBold + TextLightRed,
		Warn:  TextBold + TextYellow,
		Info:  TextGreen,
		Debug: TextGray,
	}
	defaultLog = New()
)

func New() *ColorLog {
	return &ColorLog{
		cfg:  &defaultDarkConfig,
		errW: os.Stderr,
		outW: os.Stdout,
	}
}

func (c *ColorLog) WithErr(iow io.Writer) *ColorLog {
	if iow == nil {
		c.errW = os.Stderr
	} else {
		c.errW = iow
	}
	return c
}

func (c *ColorLog) WithOut(iow io.Writer) *ColorLog {
	if iow == nil {
		c.errW = os.Stderr
	} else {
		c.outW = iow
	}
	return c
}

func (c *ColorLog) WithConfig(cfg *Config) *ColorLog {
	if cfg != nil {
		c.cfg = cfg
	}
	return c
}

func (c *ColorLog) WithDefaultDarkConfig() *ColorLog {
	c.cfg = &defaultDarkConfig
	return c
}

func (c *ColorLog) WithDefaultLightConfig() *ColorLog {
	c.cfg = &defaultLightConfig
	return c
}

func (c *ColorLog) l(isOut bool, text string, styles ...string) {
	t := time.Now()
	d := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	sb := strings.Builder{}

	if d != c.date {
		c.date = d
		sb.WriteString(c.cfg.Date + "\n")
		sb.WriteString(c.date.Format(time.DateOnly))
		sb.WriteString(TextReset + "\n")
	}

	sb.WriteString(c.cfg.Time)
	sb.WriteString(t.Format(time.TimeOnly))
	sb.WriteString(TextReset)
	for _, style := range styles {
		sb.WriteString(style)
	}
	sb.WriteString(" ")
	sb.WriteString(text)
	sb.WriteString(TextReset)
	sb.WriteString("\n")
	s := sb.String()
	if isOut {
		_, err := c.outW.Write([]byte(s))
		if err != nil {
			println(err.Error())
		}
	} else {
		_, err := c.errW.Write([]byte(s))
		if err != nil {
			println(err.Error())
		}
	}
}

func (c *ColorLog) Fatal(text string, args ...any) {
	c.l(false, fmt.Sprintf(text, args...), c.cfg.Fatal)
}

func (c *ColorLog) Error(text string, args ...any) {
	c.l(false, fmt.Sprintf(text, args...), c.cfg.Error)
}

func (c *ColorLog) Warn(text string, args ...any) {
	c.l(true, fmt.Sprintf(text, args...), c.cfg.Warn)
}

func (c *ColorLog) Info(text string, args ...any) {
	c.l(true, fmt.Sprintf(text, args...), c.cfg.Info)
}

func (c *ColorLog) Debug(text string, args ...any) {
	c.l(true, fmt.Sprintf(text, args...), c.cfg.Debug)
}

func Fatal(text string, args ...any) {
	defaultLog.Fatal(text, args...)
}

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
