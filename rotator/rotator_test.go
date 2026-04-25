package rotator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenOrCreate(t *testing.T) {
	t.Parallel()

	prefix := "testlog_open_or_create_"

	delLogs(prefix)
	defer delLogs(prefix)

	_, err := NewBuilder().WithFilePrefix(prefix).Build().Write([]byte("log1"))
	require.NoError(t, err)

	_, err = NewBuilder().WithFilePrefix(prefix).Build().Write([]byte("log2"))
	require.NoError(t, err)

	logFiles := getLogFiles(prefix)
	require.Len(t, logFiles, 1)

	r, err := os.Open(logFiles[0])
	require.NoError(t, err)
	defer func(r *os.File) {
		_ = r.Close()
	}(r)

	content, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, "log1log2", string(content))
}

func TestRotatorStrategySize(t *testing.T) {
	t.Parallel()

	prefix := "testlog_size_"

	delLogs(prefix)
	defer delLogs(prefix)

	require.NotPanics(t, func() {
		r := NewBuilder().WithStrategy(StrategySize).
			WithCount(3).
			WithSize(125).
			WithFilePrefix(prefix).
			WithCheckPeriod(time.Millisecond).
			Build()
		wg := &sync.WaitGroup{}
		wg.Add(10)
		var err error
		for g := 0; g < 10; g++ {
			go func() {
				for i := 0; i < 10000; i++ {
					_, e := r.Write([]byte(fmt.Sprintf("%d\n", i)))
					err = errors.Join(err, e)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		assert.NoError(t, err)

		logFiles := getLogFiles(prefix)
		sort.Strings(logFiles)

		require.Equal(t, len(logFiles), 3)
		assert.Equal(t, logFiles, []string{prefix + "00.log", prefix + "01.log", prefix + "02.log"})
	})
}

func TestRotatorStrategyPeriod(t *testing.T) {
	t.Parallel()

	prefix := "testlog_time_"

	delLogs(prefix)
	defer delLogs(prefix)

	require.NotPanics(t, func() {
		r := NewBuilder().WithStrategy(StrategyPeriod).
			WithCount(20).
			WithPeriod(time.Millisecond * 100).
			WithFilePrefix(prefix).
			WithCheckPeriod(time.Millisecond).
			Build()
		wg := &sync.WaitGroup{}
		wg.Add(10)
		var err error
		for g := 0; g < 10; g++ {
			go func() {
				for i := 0; i < 10; i++ {
					_, e := r.Write([]byte(fmt.Sprintf("%d\n", i)))
					err = errors.Join(err, e)
					time.Sleep(time.Millisecond * 45)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		assert.NoError(t, err)

		logFiles := getLogFiles(prefix)
		sort.Strings(logFiles)

		require.Equal(t, 4, len(logFiles))
		assert.Equal(t, []string{prefix + "00.log", prefix + "01.log", prefix + "02.log", prefix + "03.log"}, logFiles)
	})
}

func delLogs(prefix string) {
	wd, _ := os.Getwd()
	entries, _ := os.ReadDir(wd)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			_ = os.Remove(name)
		}
	}
}

func getLogFiles(prefix string) []string {
	wd, _ := os.Getwd()
	entries, _ := os.ReadDir(wd)
	files := make([]string, 0)
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, prefix) {
			files = append(files, name)
			println(name)
		}
	}
	return files
}
