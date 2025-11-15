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
			fmt.Printf("** plantuml file %s exists but is empty - overwrite? (Y/N): ", pth)
			input := ConsolePrompt()
			switch input {
			case "n":
				return fmt.Errorf("user declined to overwrite %s, exiting", pth)
			}
		} else {
			fmt.Printf("** plantuml file %s exists with %d bytes of content - overwrite? (Y/N): ", pth, fsize)
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
	fmt.Printf("wrote %d bytes to %s\n", n, pth)
	return nil
}

type PumlWriter interface {
	Out() []byte
}

type UmlClass struct {
	Title string
	r     *rgx.Rgx
}

func NewUmlClass(title string, r *rgx.Rgx) *UmlClass {
	return &UmlClass{Title: title, r: r}
}

func (d *UmlClass) Out() []byte {
	return fmt.Appendf(nil, "@startuml gopuml_class_dgm\ntitle %s\n%s\n@enduml",
		d.Title, d.BuildDiagram())
}

func (d *UmlClass) BuildDiagram() string {
	var diagramStr string

	// iterate through directories
	for pkg, file := range d.r.PkgMap {
		// create outer plantuml package to represent each directory
		pkgUmlStr := fmt.Sprintf("package %s {\n", pkg)

		// build package for each file in directory
		for fname, frgx := range file {
			// create inner package to represent file
			pkgUmlStr += (FileAsUMLPkg(fname, frgx) + "}\n")
		}
		// close directory package
		diagramStr += fmt.Sprintf("%s}\n", pkgUmlStr)
	}
	return diagramStr
}

// return a string with plantuml syntax that creates a package per file to go in UML class diagram
func FileAsUMLPkg(fname string, frgx *rgx.RgxFile) string {
	fileUmlStr := fmt.Sprintf("\tpackage %s {\n", strings.TrimSuffix(filepath.Base(fname), ".go"))
	fileUmlStr += fmt.Sprintf("%s\t\t}\n", PkgLevelFuncs(fname, frgx.Funcs))
	return fmt.Sprintf("%s\t%s", fileUmlStr, StructsInFile(frgx.Structs, frgx.Methods))
}

// return string of uml class holding all funcs that aren't methods
// funcs in a file that aren't preceeded with (s *SomeStruct)
func PkgLevelFuncs(fname string, funcs []*rgx.RgxFunc) string {
	fileFuncsStr := fmt.Sprintf("\t\tclass %s_funs {\n", strings.TrimSuffix(filepath.Base(fname), ".go"))
	for _, fn := range funcs {
		fileFuncsStr += fmt.Sprintf("\t\t\t+ %s() %s\n", fn.Name, fn.Rtn)
	}
	return fileFuncsStr
}

// return string with several uml classes, to be appended inside a package string
func StructsInFile(structs []*rgx.RgxStruct, methods []*rgx.RgxFunc) string {
	var structsStr string
	for _, s := range structs {
		structStr := fmt.Sprintf("\tclass %s {\n", s.Name)
		// append each field in the struct
		for _, fld := range s.Fields {
			structStr += fmt.Sprintf("\t\t\t+ %s %s\n", fld.Name, fld.DType)
		}
		// append each method belonging to the struct (with or without * prefix)
		for _, m := range methods {
			if strings.TrimPrefix(m.BelongsTo, "*") == s.Name {
				structStr += fmt.Sprintf("\t\t\t+ %s() %s\n", m.Name, m.Rtn)
			}
		}
		// close struct class
		structsStr += fmt.Sprintf("%s\t\t}\n\t", structStr)
	}
	return structsStr
}

type UmlActivity struct {
	Title string
	// config/styling options
}

func (d *UmlActivity) Out() []byte {
	return []byte("@startuml\n@enduml")
}
