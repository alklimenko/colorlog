package colorlog_test

import (
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/alklimenko/colorlog"
	"github.com/alklimenko/colorlog/rotator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	sb strings.Builder
}

var (
	regDate = regexp.MustCompile(`\d{4}-\d\d-\d\d`)
	regTime = regexp.MustCompile(`\d\d:\d\d:\d\d`)
)

func (t *testWriter) Write(p []byte) (n int, err error) {
	if t == nil {
		return 0, errors.New("testWriter is nil")
	}
	t.sb.WriteString(string(p))
	return len(p), nil
}

func (t *testWriter) Get() string {
	if t == nil {
		return ""
	}
	return t.sb.String()
}

func (t *testWriter) Reset() {
	t.sb.Reset()
}

func Test_Mask(t *testing.T) {
	t.Parallel()

	testData := []struct {
		name  string
		masks []string
		str   string
		res   string
	}{
		{
			name:  "разные длины",
			masks: []string{"1", "22", "333", "4444", "55555", "666666"},
			str:   "1 22 333 4444 55555 666666",
			res:   "* ** *** **** **...** **...**",
		},
		{
			name:  "2 маскированных слова",
			masks: []string{"first", "second"},
			str:   "bla-bla-bla second tru-tu first the end",
			res:   "bla-bla-bla **...** tru-tu **...** the end",
		},
		{
			name:  "большое включает малое",
			masks: []string{"one", "another one"},
			str:   "bla-bla-bla another one tru-tu one the end",
			res:   "bla-bla-bla **...** tru-tu *** the end",
		},
		{
			name:  "несколько раз одно",
			masks: []string{"one"},
			str:   "bla-bla-bla one one tru-tu one the end",
			res:   "bla-bla-bla *** *** tru-tu *** the end",
		},
	}

	for _, test := range testData {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			tw := &testWriter{}
			var l = colorlog.New().WithOut(tw)
			l.Info("") // date
			require.NotPanics(t, func() {
				for _, mask := range test.masks {
					l.WithMask(mask)
				}
				for _, logF := range []func(string){l.Fatal, l.Error, l.Info, l.Warn, l.Debug} {
					tw.Reset()
					logF(test.str)
					lstr := tw.Get()
					parts := strings.SplitN(lstr, " ", 2)
					assert.Equal(t, test.res, strings.Trim(parts[1], colorlog.TextReset+"\n"))
				}
			})
		})
	}
}

func TestLogMessage(t *testing.T) {
	t.Parallel()
	tw := &testWriter{}
	var l = colorlog.New().WithOut(tw).WithDefaultDarkConfig()
	l.Fatal("Fatal message")
	l.Error("Error message")
	l.Warn("Warning message")
	l.Info("Info message")
	l.Debug("Debug message")
	logStr := tw.Get()
	expected := "\u001B[1m\u001B[96m\n2000-01-01\u001B[0m\n\u001B[96m00:00:00\u001B[0m\u001B[1m\u001B[3m\u001B[95m Fatal message\u001B[0m\n\u001B[96m00:00:00\u001B[0m\u001B[1m\u001B[91m Error message\u001B[0m\n\u001B[96m00:00:00\u001B[0m\u001B[1m\u001B[93m Warning message\u001B[0m\n\u001B[96m00:00:00\u001B[0m\u001B[92m Info message\u001B[0m\n\u001B[96m00:00:00\u001B[0m\u001B[37m Debug message\u001B[0m\n"
	logStr = regDate.ReplaceAllString(logStr, "2000-01-01")
	logStr = regTime.ReplaceAllString(logStr, "00:00:00")

	assert.Equal(t, expected, logStr)
}

func TestLogFormattedMessage(t *testing.T) {
	t.Parallel()
	tw := &testWriter{}
	var l = colorlog.New().WithOut(tw).WithDefaultLightConfig()

	l.Fatalf("%d Fatal %s message", 1, "formatting")
	l.Errorf("Error (%d) %s", 2, "message")
	l.Warnf("Warning message: %d %v", 3, struct {
		field1 string
		field2 int
	}{field1: "field1", field2: 2})
	l.Infof("%d. %s %s", 4, "Info", "message")
	l.Debugf("Debug message without formatting")

	logStr := tw.Get()
	expected := "\u001B[1m\u001B[36m\n2000-01-01\u001B[0m\n\u001B[36m00:00:00\u001B[0m\u001B[1m\u001B[3m\u001B[95m 1 Fatal formatting message\u001B[0m\n\u001B[36m00:00:00\u001B[0m\u001B[1m\u001B[91m Error (2) message\u001B[0m\n\u001B[36m00:00:00\u001B[0m\u001B[1m\u001B[33m Warning message: 3 {field1 2}\u001B[0m\n\u001B[36m00:00:00\u001B[0m\u001B[32m 4. Info message\u001B[0m\n\u001B[36m00:00:00\u001B[0m\u001B[37m Debug message without formatting\u001B[0m\n"
	logStr = regDate.ReplaceAllString(logStr, "2000-01-01")
	logStr = regTime.ReplaceAllString(logStr, "00:00:00")

	assert.Equal(t, expected, logStr)
}

func TestRotator(t *testing.T) {
	t.Parallel()
	rot := rotator.NewBuilder().
		WithStrategy(rotator.StrategySize).
		WithSize(100).
		WithCount(5).
		WithCheckPeriod(time.Millisecond).
		Build()
	log := colorlog.New().
		WithOut(rot)

	for i := 0; i < 500; i++ {
		log.Info("Info message")
		time.Sleep(time.Millisecond)
	}
}
