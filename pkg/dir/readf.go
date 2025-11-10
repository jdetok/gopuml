package dir

import (
	"os"
	"path/filepath"
	"strings"
)

func CheckDirForFType(dir, fType string, m map[string]string) (int, error) {
	count := 0
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for _, item := range dirItems {
		path := filepath.Join(dir, item.Name())
		if item.IsDir() {
			numF, err := CheckDirForFType(path, fType, m)
			if err != nil {
				return 0, err
			}
			count += numF
			continue
		}
		if strings.HasSuffix(item.Name(), fType) {
			count++
			m[dir] = item.Name()
			// fmt.Printf("%s contains file type %s\n", path, fType)
		}
	}
	return count, nil
}
