package puml

import (
	"bufio"
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

func ConsolePrompt() string {
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input))
}

func (p *Puml) WriteOutput(dir, fname string) error {
	var err error

	// get bytes with formatted plantuml source code
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

	var f *os.File
	defer f.Close()

	// join dir and fname to path, trim any whitespace
	pth := strings.TrimSpace(filepath.Join(dir, fname))

	// if path doesn't already have puml suffix, add it
	fsuf := ".puml"
	if !strings.HasSuffix(pth, fsuf) {
		pth += fsuf
	}

	// check if file exists, create if it doesn't
	if info, err := os.Stat(pth); err == nil {
		fsize := info.Size()
		if fsize == 0 {
			fmt.Printf("%s exists but is empty, overwrite? (y/n): ", pth)
			input := ConsolePrompt()
			if input != "y" {
				return fmt.Errorf("user declined to overwrite %s, exiting", pth)
			}
		} else {
			fmt.Printf("%s exists but with %d bytes, overwrite? (y/n): ", pth, fsize)
			input := ConsolePrompt()
			if input != "y" {
				return fmt.Errorf("user declined to overwrite %s, exiting", pth)
			}
		}
		f, err = os.OpenFile(pth, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err

		}
	} else {
		fmt.Printf("%s directory does not exist, creating\n", pth)
		f, err = os.Create(pth)
		if err != nil {
			return &errd.FileCreateError{Path: pth, Err: err}
		}
		fmt.Printf("file successfully created at %s\n", dir)
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
