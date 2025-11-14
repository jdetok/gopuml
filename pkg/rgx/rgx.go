package rgx

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/jdetok/gopuml/pkg/dir"
	"github.com/jdetok/gopuml/pkg/str"
)

// TODO: different pattern for methods/funcs
const (
	// keys in rgx map
	EMPTY     = "empty" // match empty line
	PKG       = "package"
	CMNT      = "comment"   // match full line comment
	FUNC      = "func"      // match functions that are not methods on a type
	MTHD      = "method"    // match methods on a type
	STRUCT    = "struct"    // match struct types
	STRUCTFLD = "structFld" // match fields of a struct type
	ENDSTMNT  = "structEnd" // match the END of a struct type - only '}'
	TYPEMAP   = "typeMap"   // match map types (map aliased to a type)

	// vals in rgx map (regex patterns)
	RGX_EMPTY      = `^\s*$`
	RGX_PKG        = `^package\s+([A-Za-z_]*)$`
	RGX_CMNT       = `^\s*([/][/|*])\s*\w*$`
	RGX_FUNC       = `^func\s*([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\(?([^{)]*|[^{\s]+)?\)?\s*\{\s*$`
	RGX_MTHD       = `^func\s+\([A-Za-z_]+\s+([^)]*)\)\s*?([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\(?([^{)]*|[^{\s]+)?\)?\s*\{\s*$`
	RGX_STRUCT     = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_STRUCT_FLD = `^\s*([A-Za-z_]\w*)\s+([*\[\]\w.{}]+)\s*([` + "`" + `].*[` + "`" + `]|\/\/.*)?\s*$`
	// RGX_STRUCT_FLD = `^\s*([A-Za-z_]\w*)\s+([*\[\]\w.{}]+(?:\s*[\[\]\w.*{}]*)?)\s*([` + "`" + `].*[` + "`" + `]|\/\/.*)?\s*$`
	RGX_ENDSTMNT = `^\s*}\s*$`
	RGX_TYPE_MAP = `^type\s+(\w+)\s+(map\[.*)$`
)

// all regex patterns compiled ahead of time into RegexReady map type
var (
	// map regex consts to their type const
	RGX_MAP = map[string]string{
		EMPTY:     RGX_EMPTY,
		PKG:       RGX_PKG,
		CMNT:      RGX_CMNT,
		FUNC:      RGX_FUNC,
		MTHD:      RGX_MTHD,
		STRUCT:    RGX_STRUCT,
		STRUCTFLD: RGX_STRUCT_FLD,
		ENDSTMNT:  RGX_ENDSTMNT,
		TYPEMAP:   RGX_TYPE_MAP,
	}
	// check lines in this specific order
	RGX_CHECK_ORDER = []string{PKG, STRUCT, TYPEMAP, FUNC, MTHD, CMNT}
)

// think for each file i need to do a struct with arrays of funcs, structs, etc
type Rgx struct {
	ready   RgxReady
	PkgMap  RgxPkgMap
	Funcs   []*RgxFunc
	Structs []*RgxStruct
}

func NewRgx() (*Rgx, error) {
	return &Rgx{
		ready:  *CompileRgx(),
		PkgMap: RgxPkgMap{},
	}, nil
}

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

// map package string to file map
type RgxPkgMap map[string]RgxFileMap

// map file name to rgx match type, each filemap is mapped to a package name
type RgxFileMap map[string]*RgxMatch

type RgxMatch struct {
	FindType string   // func, struct, struct field, etc
	MatchStr string   // string where match was found
	Groups   []string // groups of matches
}

type RgxFunc struct {
	Name      string // function name
	Params    string // params
	Rtn       string // return types
	IsMthd    bool
	BelongsTo string
}

type RgxStruct struct {
	Name   string
	Fields []*RgxStructFld
}

type RgxStructFld struct {
	Name  string // field name
	DType string // field type
	Cmnt  string
	Tag   string
}

// compile all regex patterns in the RegexPtrns type - map returned Regexp type
// to the item name
func CompileRgx() *RgxReady {
	rr := RgxReady{}
	for name, ptrn := range RGX_MAP {
		rr[name] = regexp.MustCompile(ptrn)
	}
	return &rr
}

// same idea as one with RgxReady passed, but capture package
func (r *Rgx) Parse(dm *dir.DirMap) error {
	for d, fm := range *dm {
		r.PkgMap[d] = RgxFileMap{}
		fmt.Println(d)
		for n := range fm {
			f, err := dm.OpenFile(d, n)

			if err != nil {
				return err
			}
			defer f.Close()
			if err := r.RgxParseFile(d, f); err != nil {
				return err
			}
		}
	}
	return nil
}

// use bufio scanner to iterate through each line in passed file
func (r *Rgx) RgxParseFile(dir string, f *os.File) error {
	defer f.Close()
	fmt.Printf("parsing %s...\n", f.Name())

	if r.PkgMap[dir][f.Name()] == nil {
		r.PkgMap[dir][f.Name()] = &RgxMatch{}
	}

	// var RgxP

	lineCount := 0
	// insideStruct := false // used to handle struct fields
	// scan each line of the file
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineCount++
		line := scanner.Text() // get string of current line
		rm := r.RgxParseLine(line)
		if rm == nil {
			continue
		}

		switch rm.FindType {
		case FUNC:
			r.Funcs = append(r.Funcs,
				&RgxFunc{
					Name:   rm.Groups[0],
					Params: rm.Groups[1],
					Rtn:    rm.Groups[2],
				},
			)
		case STRUCT:
			var s = &RgxStruct{
				Name: rm.Groups[0],
			} // get fields from struct (scans lines until '}' is reached)
			matches := r.RgxParseStruct(scanner, &lineCount)
			for _, m := range matches {
				fmt.Printf("%v\n", *m)
				s.Fields = append(s.Fields,
					&RgxStructFld{
						Name:  m.Groups[0],
						DType: m.Groups[1],
					},
				)
			}
			r.Structs = append(r.Structs, s)
		case MTHD:
			r.Funcs = append(r.Funcs,
				&RgxFunc{
					IsMthd:    true,
					BelongsTo: rm.Groups[0],
					Name:      rm.Groups[1],
					Params:    rm.Groups[2],
					Rtn:       rm.Groups[3],
				},
			)
		}

	}
	return nil
}

// pass line bytes and linenum, check for regex matches
// return matches and a bool signaling whether next line is within a struct
func (r *Rgx) RgxParseLine(line string) *RgxMatch {
	var m RgxMatch
	for _, key := range RGX_CHECK_ORDER {
		rgx, ok := r.ready[key]
		if !ok {
			continue
		}
		if matches := rgx.FindStringSubmatch(line); matches != nil {
			m.FindType = key
			m.MatchStr = matches[0]
			m.Groups = matches[1:]
			groupStr := str.AngleWrap(m.Groups)
			fmt.Printf("MATCH: %s\nKEY: %s | GROUPS: %s\n", m.FindType, key, groupStr)
			if m.FindType != STRUCT {
				fmt.Println()
			}
			return &m
		}
	}
	return nil
}

func (r *Rgx) RgxParseStruct(s *bufio.Scanner, lineCount *int) []*RgxMatch {
	insideStruct := true
	matches := []*RgxMatch{}
	for s.Scan() && insideStruct {
		var rm *RgxMatch
		*lineCount++
		structLine := s.Text()

		rm, insideStruct = r.RgxParseStructFld(structLine)
		if rm == nil {
			if insideStruct {
				continue
			} else {
				break
			}
		}
		matches = append(matches, rm)
	}

	// fmt.Printf("finished scanning struct from lines %d - %d\n", startLine, *lineCount)
	return matches
}

// match fields inside a struct, only called from RgxParseStruct
// the bool returned determines whether the struct has more fields or not
func (r *Rgx) RgxParseStructFld(line string) (*RgxMatch, bool) {
	// check for ; at end of struct def - send signal to end RgxParseStruct
	if r.ready[ENDSTMNT].MatchString(line) {
		return nil, false
	}
	// if not a ; check for fields | if nil return nil but true continue signal
	matches := r.ready[STRUCTFLD].FindStringSubmatch(line)
	if matches == nil {
		return nil, true
	}
	// if field found build match struct
	return &RgxMatch{
		FindType: STRUCTFLD,
		MatchStr: matches[0],
		Groups:   matches[1:],
	}, true
}
