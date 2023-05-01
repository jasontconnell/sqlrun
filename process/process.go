package process

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func GetSqlFiles(dir, priority, exclude, filter, ext string) ([]string, error) {
	pmap := make(map[string]int)

	dofilter := filter != ""

	plist := strings.Split(priority, ",")
	for i, p := range plist {
		pmap[p] = i
	}

	var elist []string
	if exclude != "" {
		elist = strings.Split(exclude, ",")
	}

	files := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		nm := info.Name()
		fext := filepath.Ext(nm)
		if fext != "."+ext { // not the file we're looking for.
			return nil
		}

		if dofilter && !strings.Contains(nm, filter) {
			return nil
		}

		if len(elist) != 0 {
			for _, exc := range elist {
				if strings.HasPrefix(nm, exc) {
					return nil
				}
			}
		}

		files = append(files, path)

		return nil
	})

	sort.Slice(files, func(i, j int) bool {
		_, fn1 := filepath.Split(files[i])
		_, fn2 := filepath.Split(files[j])
		p1 := getPrefixSort(pmap, fn1)
		p2 := getPrefixSort(pmap, fn2)
		val := p1 < p2
		if p1 == 999 && p2 == 999 {
			val = fn1 < fn2
		}
		return val
	})

	return files, nil
}

func getPrefixSort(pmap map[string]int, v string) int {
	sortVal := 999
	for prefix, i := range pmap {
		if strings.HasPrefix(v, prefix) {
			sortVal = i
		}
	}
	return sortVal
}
