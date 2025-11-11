package dir

import (
	"fmt"
	"os"

	"github.com/jdetok/gopuml/pkg/errd"
)

func (dm DirMap) ReadFile(dir, file string) ([]byte, error) {
	return os.ReadFile(dm[dir][file])
}

func (dm DirMap) ReadAll() error {
	for dir, files := range dm {
		for name, path := range files {
			b, err := dm.ReadFile(dir, name)
			if err != nil {
				return &errd.FileReadError{Path: path, Err: err}
			}
			fmt.Println(string(b))
		}
	}
	return nil
}
