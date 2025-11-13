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
	ENDSTMNT  = "structEnd" // match the END of a struct type - only '}'
	TYPEMAP   = "typeMap"   // match map types (map aliased to a type)

	// vals in rgx map (regex patterns)
	RGX_EMPTY      = `^\s*$`
	RGX_CMNT       = `^\s*([/][/|*])\s*\w*$`
	RGX_FUNC       = `^func\s+([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\s*(?:([^{]+))\s*\{\s*$`
	RGX_MTHD       = `^func\s+(?:\(([^)]*)\)\s*)([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\s*(?:([^{]+))\s*\{\s*$`
	RGX_STRUCT     = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_STRUCT_FLD = `^\s*([A-Za-z_]\w*)\s+([*\[\]\w.{}]+(?:\s*[\[\]\w.*{}]*)?)\s*([` + "`" + `].*[` + "`" + `]|\/\/.*)?\s*$`
	RGX_ENDSTMNT   = `^\s*}\s*$`
	RGX_TYPE_MAP   = `^type\s+(\w+)\s+(map\[.*)$`
)

// all regex patterns compiled ahead of time into RegexReady map type
var (
	// map regex consts to their type const
	RGX_MAP = map[string]string{
		EMPTY:     RGX_EMPTY,
		CMNT:      RGX_CMNT,
		FUNC:      RGX_FUNC,
		MTHD:      RGX_MTHD,
		STRUCT:    RGX_STRUCT,
		STRUCTFLD: RGX_STRUCT_FLD,
		ENDSTMNT:  RGX_ENDSTMNT,
		TYPEMAP:   RGX_TYPE_MAP,
	}
	// check lines in this specific order
	RGX_CHECK_ORDER = []string{STRUCT, TYPEMAP, FUNC, MTHD, CMNT}
)

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

type RgxFileMap map[string]RgxLineMap
type RgxLineMap map[int]*RgxMatch

type RgxMatch struct {
	FindType string   // func, struct, struct field, etc
	MatchStr string   // string where match was found
	Groups   []string // groups of matches
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
	fmt.Printf("parsing %s...\n", f.Name())
	scanner := bufio.NewScanner(f) // scans one line at a time

	lineCount := 0
	insideStruct := false
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		if !insideStruct {
			rm := rr.RgxParseLine(line)
			if rm == nil {
				continue
			}

			switch rm.FindType {
			case STRUCT:
				structLines := rr.RgxParseStruct(scanner, &lineCount)
				if structLines != nil {
					fmt.Println(len(structLines), "lines from struct")
				}
			}
		}
	}

	return nil
}

// pass line bytes and linenum, check for regex matches
// return matches and a bool signaling whether next line is within a struct
func (rr RgxReady) RgxParseLine(line string) *RgxMatch {
	var m RgxMatch
	for _, key := range RGX_CHECK_ORDER {
		rgx, ok := rr[key]
		if !ok {
			continue
		}
		if matches := rgx.FindStringSubmatch(line); matches != nil {
			m.FindType = key
			m.MatchStr = matches[0]
			m.Groups = matches[1:]
			groupStr := "<" + strings.Join(m.Groups, "> <") + ">"
			fmt.Printf("MATCH: %s\nKEY: %s | GROUPS: %s\n\n", m.FindType, key, groupStr)
			return &m
		}
	}
	return nil
}

func (rr RgxReady) RgxParseStruct(s *bufio.Scanner, lineCount *int) [][]string {
	insideStruct := true
	startLine := *lineCount
	var lines [][]string

	for s.Scan() && insideStruct {
		*lineCount++
		var matches []string
		structLine := s.Text()

		matches, insideStruct = rr.RgxParseStructFld(structLine)
		if matches == nil || matches[0] == EMPTY {
			continue
		}
		if matches[0] == ENDSTMNT {
			break
		}
		lines = append(lines, matches)
	}

	fmt.Printf("finished scanning struct from lines %d - %d\n", startLine, *lineCount)
	return lines
}

func (rr RgxReady) RgxParseStructFld(line string) ([]string, bool) {
	results := []string{}
	if ok := rr[ENDSTMNT].MatchString(line); ok {
		return append(results, ENDSTMNT), false
	}
	if ok := rr[EMPTY].MatchString(line); ok {
		return append(results, EMPTY), true
	}
	if results = rr[STRUCTFLD].FindStringSubmatch(line); results != nil {
		return results, true
	}
	return nil, false
}
