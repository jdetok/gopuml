package parse

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

const (
	RGX_STRUCT   = `^type\s+(.+?)\s*struct\s*{\s*$`
	RGX_TYPE_MAP = `^type\s+(.+?)\s*map\s*(.+?)$`
	RGX_FUNC     = `^func(.+?)\s*\{\s*$`
)

type RegexPtrns map[string]string
type RegexReady map[string]*regexp.Regexp

func NewRegexPtrns() *RegexPtrns {
	return &RegexPtrns{
		"func":    RGX_FUNC,
		"struct":  RGX_STRUCT,
		"typeMap": RGX_TYPE_MAP,
	}
}

func (rp *RegexPtrns) CompileRegex() (*RegexReady, error) {
	rr := RegexReady{}
	for name, ptrn := range *rp {
		rgx, err := regexp.Compile(ptrn)
		if err != nil {
			return nil, err
		}
		rr[name] = rgx
	}
	return &rr, nil
}

func (rr *RegexReady) FileParser(f *os.File) error {
	defer f.Close()
	scanner := bufio.NewScanner(f) // scans one line at a time

	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		rr.ParseLine(line, lineNum)
	}
	return nil
}

func (rr *RegexReady) ParseLine(line []byte, lineNum int) {
	for name, rgx := range *rr {
		if ok := rgx.Match(line); ok {
			fmt.Printf("%s found at %d\n", name, lineNum)
		}
	}
}
