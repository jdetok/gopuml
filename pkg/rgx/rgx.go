package rgx

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO: different pattern for methods/funcs
const (
	// keys in rgx map
	EMPTY     = "empty"     // match empty line
	CMNT      = "comment"   // match full line comment
	FUNC      = "func"      // match functions that are not methods on a type
	MTHD      = "method"    // match methods on a type
	STRUCT    = "struct"    // match struct types
	STRUCTFLD = "structFld" // match fields of a struct type
	STRUCTEND = "structEnd" // match the END of a struct type - only '}'
	TYPEMAP   = "typeMap"   // match map types (map aliased to a type)

	// vals in rgx map (regex patterns)
	RGX_EMPTY      = `^\s*$`
	RGX_CMNT       = `^\s*([/][/|*])\s*\w*$`
	RGX_FUNC       = `^func\s+(?:\(([^)]*)\)\s*)?([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\s*(?:([^{]+))\s*\{\s*$`
	RGX_MTHD       = `^func\s+(?:\(([^)]*)\)\s*)?([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\s*(?:([^{]+))\s*\{\s*$`
	RGX_STRUCT     = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_STRUCT_FLD = `^(?:[^func]\s*)([A-Za-z_]\w*)\s+([^` + "`" + `\s/]+(?:\s+[^` + "`" + `\s/]+)*)\s*(.*)$`
	RGX_STRUCT_END = `^\s*}\s*$`
	RGX_TYPE_MAP   = `^type\s+(\w+)\s+(map\[.*)$`
)

// all regex patterns compiled ahead of time into RegexReady map type
var RGX_MAP = map[string]string{
	EMPTY:     RGX_EMPTY,
	CMNT:      RGX_CMNT,
	FUNC:      RGX_FUNC,
	MTHD:      RGX_MTHD,
	STRUCT:    RGX_STRUCT,
	STRUCTFLD: RGX_STRUCT_FLD,
	STRUCTEND: RGX_STRUCT_END,
	TYPEMAP:   RGX_TYPE_MAP,
}

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

type RgxFileMap map[string]RgxLineMap
type RgxLineMap map[int]*RgxMatch

type RgxMatch struct {
	FindType string        // func, struct, struct field, etc
	RawStr   string        // string where match was found
	Groups   []RgxMatchGrp // groups of matches
}

// can be one or the other
type RgxMatchGrp struct {
	Fn RgxFuncGrp   // function struct
	St RgxStructGrp // struct struct
}

type RgxFuncGrp struct {
	FnName   string   // function name
	FnSign   string   // function signature
	FnParams []string // params
	FnRtn    []string // return types
}

type RgxStructGrp struct {
	Fields []RgxStructFldGrp
}

type RgxStructFldGrp struct {
	Name  string // field name
	DType string // field type
}

// compile all regex patterns in the RegexPtrns type - map returned Regexp type
// to the item name
func CompileRgx() (*RgxReady, error) {
	rr := RgxReady{}
	for name, ptrn := range RGX_MAP {
		rgx, err := regexp.Compile(ptrn)
		if err != nil {
			return nil, err
		}
		rr[name] = rgx
	}
	return &rr, nil
}

// use bufio scanner to iterate through each line in passed file
func (rr RgxReady) RgxParseFile(f *os.File) error {
	defer f.Close()
	scanner := bufio.NewScanner(f) // scans one line at a time

	lineCount := 0
	// todo - set when inside struct (after struct is found)
	// USE SWITCH - just switch no var
	insideStruct := false
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		if !insideStruct {
			res := rr.RgxParseLine(line)
			if res == nil {
				continue
			}
			fmt.Println(res[0])
			if strings.Contains(res[0], "struct") {
				structLines := rr.RgxParseStruct(scanner, &lineCount)
				if structLines != nil {
					fmt.Println(len(structLines), "lines from struct")
					for _, line := range structLines {
						fmt.Println(line[0])
					}
				}
			}
		}
	}

	return nil
}

// pass line bytes and linenum, check for regex matches
// return matches and a bool signaling whether next line is within a struct
func (rr RgxReady) RgxParseLine(line string) []string {
	for _, rgx := range rr {
		matches := rgx.FindStringSubmatch(line)
		if matches == nil {
			return nil
		}
		return matches
	}
	return nil
}

func (rr RgxReady) RgxParseStruct(s *bufio.Scanner, lineCount *int) [][]string {
	insideStruct := true
	var lines [][]string
	fmt.Println("started struct scan at line", *lineCount)
	for s.Scan() && insideStruct {
		*lineCount++
		var matches []string
		structLine := s.Text()
		matches, insideStruct = rr.RgxParseStructFld(structLine)
		if matches == nil {
			continue // might want to change this
		}
		if matches[0] == EMPTY {
			continue
		}
		if matches[0] == STRUCTEND {
			return lines
		}
		lines = append(lines, matches)
	}
	fmt.Println("finished struct scan at line", *lineCount)
	return lines
}

func (rr RgxReady) RgxParseStructFld(line string) ([]string, bool) {
	results := []string{}
	if ok := rr[STRUCTEND].MatchString(line); ok {
		return append(results, STRUCTEND), false
	}
	if ok := rr[EMPTY].MatchString(line); ok {
		return append(results, EMPTY), true
	}
	if results = rr[STRUCTFLD].FindStringSubmatch(line); results != nil {
		return results, true
	}
	return nil, false
}
