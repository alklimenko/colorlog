package rotator

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func getNext(dirName, namePrefix, ext string, count int) (string, error) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	dir, err := os.Open(dirName)
	if err != nil {
		return "", err
	}
	defer func(dir *os.File) {
		_ = dir.Close()
	}(dir)

	entries, err := dir.ReadDir(-1)
	if err != nil {
		return "", err
	}
	logfilesMap := map[int]string{}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasPrefix(name, namePrefix) || !strings.HasSuffix(name, ext) {
			continue
		}

		numStr := strings.TrimSuffix(strings.TrimPrefix(name, namePrefix), ext)
		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}
		logfilesMap[num] = name
	}

	var keys []int
	for k := range logfilesMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	idx := make(map[int]int, len(keys))
	for i, k := range keys {
		idx[i] = k
	}

	logfiles := make([]string, len(logfilesMap))
	for i := range keys {
		logfiles[i] = logfilesMap[idx[i]]
	}

	var i int
	for i = len(logfiles) - 1; i >= count; i-- {
		_ = os.Remove(filepath.Join(dirName, logfiles[i]))
	}

	for ; i >= 0; i-- {
		newName := getFilename(dirName, namePrefix, ext, i+1)
		_ = os.Rename(filepath.Join(dirName, logfiles[i]), newName)
	}

	return getFilename(dirName, namePrefix, ext, 0), nil
}

func getFilename(dirName, namePrefix, ext string, num int) string {
	return filepath.Join(dirName, fmt.Sprintf("%s%02d%s", namePrefix, num, ext))
}
