package rgx

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// REGEX PATTERNS
const (
	RGX_STRUCT     = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_STRUCT_FLD = `^\s*([A-Za-z_]\w*)\s+([^` + "`" + `\s/]+(?:\s+[^` + "`" + `\s/]+)*)\s*(.*)$`
	RGX_STRUCT_END = `^\s*}\s*$`
	RGX_TYPE_MAP   = `^type\s+(\w+)\s+(map\[.*)$`
	RGX_FUNC       = `^func\s+(?:\(([^)]*)\)\s*)?([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\s*(?:([^{]+))\s*\{\s*$`
)

// map name of item (func, struct, etc) to its regex pattern const
// all regex patterns compiled ahead of time into RegexReady map type
type RgxPtrns map[string]string

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

// map regex patterns to item names here, then call rp.CompileRegex to get regexp
func NewRgxPtrns() *RgxPtrns {
	return &RgxPtrns{
		"func":      RGX_FUNC,
		"struct":    RGX_STRUCT,
		"structFld": RGX_STRUCT_FLD,
		"structEnd": RGX_STRUCT_END,
		"typeMap":   RGX_TYPE_MAP,
	}
}

// compile all regex patterns in the RegexPtrns type - map returned Regexp type
// to the item name
func (rp *RgxPtrns) CompileRgx() (*RgxReady, error) {
	rr := RgxReady{}
	for name, ptrn := range *rp {
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

	lineNum := 0
	// todo - set when inside struct (after struct is found)
	// USE SWITCH - just switch no var
	// insideStruct := false
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		res := rr.RgxParseLine(line, lineNum)
		if res == "" {
			continue
		}
		if res == "struct" {

			structLines := rr.RgxParseStruct(scanner)
			for _, sline := range structLines {
				fmt.Printf("struct field: %s\n", sline)
			}
			// for scanner.Scan() {
			// 	structLine := scanner.Bytes()
			// 	if ok := rr["structEnd"].Match(structLine); ok {
			// 		break
			// 	}
			// 	// structLines = append(structLines, structLine)
			// }
			// for _, l := range structLines {
			// 	fmt.Println(string(l))
			// }
		}
	}
	return nil
}

// pass line bytes and linenum, check for regex matches
func (rr *RgxReady) RgxParseLine(line []byte, lineNum int) string {
	for name, rgx := range *rr {
		if ok := rgx.Match(line); ok {
			fmt.Printf("%s found at %d\n", name, lineNum)
			return name
			// if name == "struct" {

			// }
		}
		// matches := rgx.FindStringSubmatch(string(line))
		// if len(matches) > 0 {
		// 	for i, m := range matches {
		// 		fmt.Printf("group %d: %s\n", i+1, m)
		// 	}
		// }
	}
	return ""
}

func (rr RgxReady) RgxParseStruct(scanner *bufio.Scanner) []string {
	var lines []string
	for scanner.Scan() {
		structLine := scanner.Bytes()
		if ok := rr["structEnd"].Match(structLine); ok {
			return lines
		}
		line := scanner.Bytes()
		if ok := rr["structFld"].Match(line); ok {
			lines = append(lines, string(line))
		}
	}
	return lines
}
