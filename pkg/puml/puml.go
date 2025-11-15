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
	"github.com/jdetok/gopuml/pkg/rgx"
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
		// FILE EXISTS: ask user whether to overwrite or exit
		fsize := info.Size()
		if fsize == 0 {
			fmt.Printf("** plantuml (.puml) file %s exists but is empty - overwrite? (Y/N): ", pth)
			input := ConsolePrompt()
			switch input {
			case "n":
				return fmt.Errorf("user declined to overwrite %s, exiting", pth)
			}
		} else {
			fmt.Printf("** plantuml (.puml) file %s exists with %d bytes of content - overwrite? (Y/N): ", pth, fsize)
			input := ConsolePrompt()
			switch input {
			case "n":
				return fmt.Errorf("user declined to overwrite %s, exiting", pth)
			}
		}
		// open EXISTING file only after confirming whether overwrite is ok
		f, err = os.OpenFile(pth, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return &errd.FileReadError{Path: pth, Err: err}
		}
	} else {
		// FILE DOES NOT EXIST: CREATE AT THE GIVEN PATH
		f, err = os.Create(pth)
		if err != nil {
			return &errd.FileCreateError{Path: pth, Err: err}
		}
		fmt.Printf("plantuml file successfully created at %s\n", pth)
	}

	// get bytes with formatted plantuml source code (PumlWriter interface)
	b := p.Dgm.Out()

	// write bytes from PumlWriter implementation to file via io.Writer interface
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
	r     *rgx.Rgx
	// dm    *dir.DirMap
	// config/styling options
}

func NewUmlClass(title string, r *rgx.Rgx) *UmlClass {
	return &UmlClass{Title: title, r: r}
}

func (d *UmlClass) Out() []byte {
	return fmt.Appendf(nil, "@startuml %s\n%s\n@enduml", d.Title, d.BuildDiagram())
}

func (d *UmlClass) BuildDiagram() string {
	indent := "\t"
	dblIndent := "\t\t"
	tplIndent := "\t\t\t"
	var diagramStr string
	for pkg, file := range d.r.PkgMap {
		// create outer plantuml package to represent each directory
		pkgUmlStr := fmt.Sprintf("package %s {\n", pkg)
		for fname, frgx := range file {

			// create inner package to represent file
			fileUmlStr := fmt.Sprintf("%spackage %s {\n", indent, strings.TrimSuffix(filepath.Base(fname), ".go"))

			// create plantuml class to represent functions in file not belonging to a struct
			fileFuncsStr := fmt.Sprintf("%sclass %s_funs {\n", dblIndent, strings.TrimSuffix(filepath.Base(fname), ".go"))
			for _, fn := range frgx.Funcs {
				fileFuncsStr += fmt.Sprintf("%s%s() %s\n", tplIndent, fn.Name, fn.Rtn)
			}
			fileUmlStr += fmt.Sprintf("%s%s}\n", fileFuncsStr, dblIndent)

			// create plantuml class for each struct in file
			for _, s := range frgx.Structs {
				structStr := fmt.Sprintf("%sclass %s {\n", dblIndent, s.Name)
				// append each field in the struct
				for _, fld := range s.Fields {
					structStr += fmt.Sprintf("%s%s %s\n", tplIndent, fld.Name, fld.DType)
				}
				// append each method belonging to the struct (with or without * prefix)
				for _, m := range frgx.Methods {
					if strings.TrimPrefix(m.BelongsTo, "*") == s.Name {
						structStr += fmt.Sprintf("%s%s() %s\n", tplIndent, m.Name, m.Rtn)
					}
				}
				// close struct class
				fileUmlStr += fmt.Sprintf("%s%s}\n", structStr, dblIndent)
			}
			// close file package
			pkgUmlStr += fmt.Sprintf("%s%s}\n", fileUmlStr, indent)
		}
		// close directory package
		diagramStr += fmt.Sprintf("%s}\n", pkgUmlStr)
	}
	return diagramStr
}

type UmlActivity struct {
	Title string
	// config/styling options
}

func (d *UmlActivity) Out() []byte {
	return []byte("@startuml\n@enduml")
}
