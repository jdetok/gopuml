package dir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func CheckDirForFType(dir, fType string) (int, error) {
	count := 0
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	for _, item := range dirItems {
		path := filepath.Join(dir, item.Name())
		if item.IsDir() {
			numF, err := CheckDirForFType(path, fType)
			if err != nil {
				return 0, err
			}
			count += numF
			continue
		}
		if strings.HasSuffix(item.Name(), fType) {
			count++
			fmt.Printf("%s contains file type %s\n", path, fType)
		}
	}

	return count, nil
}
