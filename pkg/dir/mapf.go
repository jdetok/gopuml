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
type FileMap map[string]string

// return DirFiles (maps files to their dir) mapped recursively
// .git directory is excluded | only files ending in ftyp are mapped
func MapFiles(rootDir, ftyp string, excludeDirs []string) (*DirMap, error) {
	// make sure rootDir exists
	_, err := os.Open(rootDir)
	if err != nil {
		return nil, &errd.FileOpenError{Path: rootDir, Err: err}
	}

	// init empty DirMap to return
	dm := DirMap{}
	// begin recursive call(s)
	if err := dm.MapRecur(rootDir, ftyp, excludeDirs); err != nil {
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
func (dm DirMap) MapRecur(dir, ftyp string, excludeDirs []string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	// iterate through entries in directory
	for _, e := range entries {
		path := filepath.Join(dir, e.Name()) // join file to root path
		if e.IsDir() && !exclude(e.Name(), excludeDirs) {
			if err := dm.MapRecur(path, ftyp, excludeDirs); err != nil { // recursive call
				return &errd.FileRecursionError{Path: path, Ftyp: ftyp, Err: err}
			}
			continue // after finishing recursive call
		} else {
			// check if item in dir ends in ftyp
			if strings.HasSuffix(e.Name(), ftyp) {
				fAbs, err := filepath.Abs(path)
				if err != nil {
					return &errd.FileAbsError{Path: path, Err: err}
				}
				// ensure map exists
				if dm[dir] == nil {
					dm[dir] = FileMap{}
				}
				// map *os.File returned from os.Open() to file name
				dm[dir][e.Name()] = fAbs
			}
		}
	}
	return nil
}

func exclude(checkDir string, excludeDirs []string) bool {
	for _, dir := range excludeDirs {
		if checkDir == dir {
			return true
		}
	}
	return false
}
