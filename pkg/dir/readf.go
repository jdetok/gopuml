package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jdetok/gopuml/pkg/errd"
)

// Files type is an alias for map[string]string - map file names to dirs
type DirFiles map[string]map[string]*os.File

// return DirFiles (maps files to their dir) mapped recursively
// .git directory is excluded | only files ending in ftyp are mapped
func MapFiles(rootDir, ftyp string) (*DirFiles, error) {
	df := DirFiles{}
	df[rootDir] = map[string]*os.File{}
	if err := df.MapRecur(rootDir, ftyp); err != nil {
		return nil, &errd.FileRecursionError{Path: rootDir, Ftyp: ftyp, Err: err}
	}
	return &df, nil
}

// recursive function to map matching files to their appropriate directory
// iterate through entries in directory - for each item, if it's a
// directory (excluding .git), a recursive call is made. otherwise, check
// whether file ends in ftyp (HasSuffix). if the suffix matches the ftyp,
// os.Open() is called on the file and the *os.File is mapped to the file name
// a DirFiles type (alias for a map) must be declared and initialized first
func (df DirFiles) MapRecur(dir, ftyp string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// iterate through entries in directory
	for _, e := range entries {
		path := filepath.Join(dir, e.Name()) // join file to root path
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".git") {
			if err := df.MapRecur(path, ftyp); err != nil { // recursive call
				return &errd.FileRecursionError{Path: path, Ftyp: ftyp, Err: err}
			}
			continue // after finishing recursive call
		} else {
			// check if item in dir ends in ftyp
			if strings.HasSuffix(e.Name(), ftyp) {
				f, err := os.Open(path)
				if err != nil {
					return &errd.FileOpenError{FName: path, Err: err}
				}
				if df[dir] == nil {
					df[dir] = map[string]*os.File{}
				}
				df[dir][e.Name()] = f
			}
		}
	}
	return nil
}
