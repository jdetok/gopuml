package rgx

import "regexp"

const (
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

// compile all regex patterns in the RegexPtrns type - map returned Regexp type
// to the item name
func CompileRgx() *RgxReady {
	rr := RgxReady{}
	for name, ptrn := range RGX_MAP {
		rr[name] = regexp.MustCompile(ptrn)
	}
	return &rr
}
