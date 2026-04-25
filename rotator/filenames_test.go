package rotator

import (
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetNext(t *testing.T) {
	t.Parallel()
	prefix := "log_getnext_"
	ext := ".txt"
	delLogs(prefix)
	defer delLogs(prefix)

	// unique, non ordered rand ints
	nums := make([]int, 10)
	for i := range nums {
		nums[i] = rand.Intn(100)*10 + i
	}

	wd, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		filename := getFilename(wd, prefix, ext, nums[i])
		file, _ := os.Create(filename)
		_ = file.Close()
	}
	fn := filepath.Join(wd, prefix+"00"+ext)
	for range 7 {
		f, e := getNext(wd, prefix, ext, 3)
		require.NoError(t, e)
		assert.Equal(t, fn, f)
	}
	files := getLogFiles(prefix)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, []string{prefix + "01" + ext, prefix + "02" + ext}, files)
}

func Test_getFilename(t *testing.T) {
	t.Parallel()
	filename := getFilename("aaa", "log_getfilename_", ".logfile", 14)
	assert.Equal(t, "aaa/log_getfilename_14.logfile", filename)
}
