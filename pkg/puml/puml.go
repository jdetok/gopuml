package puml

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jdetok/gopuml/pkg/errd"
)

type Puml struct {
	Dgm PumlWriter // different diagram structs
}

func (p *Puml) WriteOutput(dir, fname string) error {
	b := p.Dgm.Out()

	// check that directory exists
	if _, err := os.Stat(dir); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("%s directory does not exist, creating\n", dir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return &errd.FileCreateError{Path: dir, Err: err}
			}
			fmt.Printf("directory created at %s\n", dir)
		} else {
			return &errd.FileReadError{Path: dir, Err: err}
		}
	}

	var err error
	var fsuf string = ".puml"
	var f *os.File
	defer f.Close()

	// build full file path
	pth := strings.TrimSpace(filepath.Join(dir, fname))
	if !strings.HasSuffix(pth, fsuf) {
		pth += fsuf
	}

	if f, err = os.OpenFile(pth, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		fmt.Printf("%s directory does not exist, creating\n", pth)
		f, err = os.Create(pth)
		if err != nil {
			return &errd.FileCreateError{Path: pth, Err: err}
		}
		fmt.Printf("file created at %s\n", dir)
	}

	n, err := f.Write(b)
	if err != nil {
		return &errd.FileWriteError{Path: fname, Err: err}
	}
	fmt.Printf("wrote %d bytes to %s\n", n, fname)
	return nil
}

type PumlWriter interface {
	Out() []byte
}

type UmlClass struct {
	Title string
	// config/styling options
}

func (d *UmlClass) Out() []byte {
	return fmt.Appendf(nil, "@startuml %s\nclass Test {\n\ttest string\n}\n@enduml", d.Title)
}

type UmlActivity struct {
	Title string
	// config/styling options
}

func (d *UmlActivity) Out() []byte {
	return []byte("@startuml\n@enduml")
}
