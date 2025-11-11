package dir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jdetok/gopuml/pkg/errd"
)

// DirMap is an alias for a nested map - directories mapped to file names,
// which are mapped to *os.File (FileMap)
type DirMap map[string]FileMap

// FileMap maps a file name to a *os.File - FileMaps are mapped to DirMap
type FileMap map[string]*os.File

// return DirFiles (maps files to their dir) mapped recursively
// .git directory is excluded | only files ending in ftyp are mapped
func MapFiles(rootDir, ftyp string) (*DirMap, error) {
	dm := DirMap{}
	if err := dm.MapRecur(rootDir, ftyp); err != nil {
		return nil, &errd.FileRecursionError{Path: rootDir, Ftyp: ftyp, Err: err}
	}
	return &dm, nil
}

// recursive function to map matching files to their appropriate directory
// iterate through entries in directory - for each item, if it's a
// directory (excluding .git), a recursive call is made. otherwise, check
// whether file ends in ftyp (HasSuffix). if the suffix matches the ftyp,
// os.Open() is called on the file and the *os.File is mapped to the file name
// a DirFiles type (alias for a map) must be declared and initialized first
func (dm DirMap) MapRecur(dir, ftyp string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// iterate through entries in directory
	for _, e := range entries {
		path := filepath.Join(dir, e.Name()) // join file to root path
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".git") {
			if err := dm.MapRecur(path, ftyp); err != nil { // recursive call
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
				if dm[dir] == nil {
					dm[dir] = FileMap{}
				}
				dm[dir][e.Name()] = f
			}
		}
	}
	return nil
}
