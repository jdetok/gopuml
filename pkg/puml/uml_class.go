package puml

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/jdetok/gopuml/pkg/rgx"
)

// BUILD UML CLASS DIAGRAM STRING

func (d *UmlClass) Build() []byte {
	diagramStr := fmt.Sprintf("@startuml %s\n", d.Title)
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
	return []byte((diagramStr + "\n@enduml"))
}

// return a string with plantuml syntax that creates a package per file to go in UML class diagram
func FileAsUMLPkg(fname string, frgx *rgx.RgxFile) string {
	fileUmlStr := fmt.Sprintf("\tpackage %s {\n", strings.TrimSuffix(filepath.Base(fname), ".go"))
	if len(frgx.Funcs) > 0 {
		fileUmlStr += fmt.Sprintf("%s\t\t}\n", PkgLevelFuncs(fname, frgx.Funcs))
	}
	return fmt.Sprintf("%s\t%s", fileUmlStr, StructsInFile(frgx.Structs, frgx.Methods))
}

// return string of uml class holding all funcs that aren't methods
// funcs in a file that aren't preceeded with (s *SomeStruct)
func PkgLevelFuncs(fname string, funcs []*rgx.RgxFunc) string {
	fileFuncsStr := fmt.Sprintf("\t\tclass %s_funcs {\n", strings.TrimSuffix(filepath.Base(fname), ".go"))
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
