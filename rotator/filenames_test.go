package rotator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_GetNext(t *testing.T) {
	t.Parallel()
	wd, _ := os.Getwd()
	wd = filepath.Join(wd, "logs")
	f, e := getNext(wd, "log_", "log", 3)
	if e != nil {
		t.Error(e)
	}
	println(f)
}

func Test2(t *testing.T) {
	t.Parallel()
	println(fmt.Sprintf("%04d", 23))
}
