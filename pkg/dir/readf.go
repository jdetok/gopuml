package dir

import (
	"fmt"
	"os"

	"github.com/jdetok/gopuml/pkg/errd"
)

// read a file in the DirMap - os.ReadFile defer file closure
func (dm DirMap) ReadFile(dir, file string) ([]byte, error) {
	return os.ReadFile(dm[dir][file])
}

func (dm DirMap) OpenFile(dir, file string) (*os.File, error) {
	return os.Open(dm[dir][file])
}

// read the contents of each file in DirMap
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
