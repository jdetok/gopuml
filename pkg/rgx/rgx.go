package rgx

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

// REGEX PATTERNS
const (
	RGX_STRUCT   = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_TYPE_MAP = `^type\s+(.+?)\s*map\s*(.+?)$`
	RGX_FUNC     = `^func(.+?)\s*\{\s*$`
)

// map name of item (func, struct, etc) to its regex pattern const
// all regex patterns compiled ahead of time into RegexReady map type
type RgxPtrns map[string]string

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

// map regex patterns to item names here, then call rp.CompileRegex to get regexp
func NewRgxPtrns() *RgxPtrns {
	return &RgxPtrns{
		"func":    RGX_FUNC,
		"struct":  RGX_STRUCT,
		"typeMap": RGX_TYPE_MAP,
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
func (rr *RgxReady) RgxParseFile(f *os.File) error {
	defer f.Close()
	scanner := bufio.NewScanner(f) // scans one line at a time

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		rr.RgxParseLine(line, lineNum)
	}
	return nil
}

// pass line bytes and linenum, check for regex matches
func (rr *RgxReady) RgxParseLine(line []byte, lineNum int) {
	for name, rgx := range *rr {
		if ok := rgx.Match(line); ok {
			fmt.Printf("%s found at %d\n", name, lineNum)
		}
	}
}
