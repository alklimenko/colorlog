package colorlog_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/alklimenko/colorlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	sb strings.Builder
}

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
			var l = colorlog.New().WithOut(tw).WithErr(tw)
			l.Info("") // date
			require.NotPanics(t, func() {
				for _, mask := range test.masks {
					l.WithMask(mask)
				}
				for _, logF := range []func(string, ...any){l.Fatal, l.Error, l.Info, l.Warn, l.Debug} {
					tw.Reset()
					logF("%s", test.str)
					lstr := tw.Get()
					parts := strings.SplitN(lstr, " ", 2)
					assert.Equal(t, test.res, strings.Trim(parts[1], colorlog.TextReset+"\n"))
				}
			})
		})
	}
}

func TestLog2(t *testing.T) {
	colorlog.Fatal("Fatal message")
	colorlog.Error("Error message")
	colorlog.Warn("Warning message")
	colorlog.Info("Info message")
	colorlog.Debug("Debug message")
}
