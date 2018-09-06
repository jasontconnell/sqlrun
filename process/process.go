package process

import (
	"os"
	"path/filepath"
	"strings"
	"sort"
)

func GetSqlFiles(dir, priority, ext string) ([]string, error) {
	pmap := make(map[string]int)
	plist := strings.Split(priority, ",")
	for i, p := range plist {
		pmap[p] = i
	}

	files := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		nm := info.Name()
		if info.IsDir() {
			return nil
		}

		fext := filepath.Ext(nm)
		if fext != "."+ext { // not the file we're looking for.
			return nil
		}

		files = append(files, path)

		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		return getPrefixSort(pmap, files[i]) < getPrefixSort(pmap, files[j])
	})

	return files, nil
}

func getPrefixSort(pmap map[string]int, v string) int {
	_, fn := filepath.Split(v)
	sortVal := 999
	for prefix, i:= range pmap {
		if strings.HasPrefix(fn, prefix) {
			sortVal = i
		}
	}
	return sortVal
}