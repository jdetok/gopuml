package rgx

import "regexp"

const (
	EMPTY           = "empty" // match empty line
	RGX_EMPTY       = `^\s*$`
	PKG             = "package" // match package
	RGX_PKG         = `^package\s+([A-Za-z_]*)$`
	IMP             = "import"
	RGX_IMP         = `^import\s+([("])([^"\s]+)?"?` // check if ( or "
	IMPKG           = "imported_pkg"
	RGX_IMPKG       = `^\s+"(\S+)"`
	CMNT            = "comment" // match full line comment
	RGX_CMNT        = `^\s*([/][/|*])\s*\w*$`
	FUNC            = "func" // match functions that are not methods on a type
	RGX_FUNC        = `^func\s*([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\(?([^{)]*|[^{\s]+)?\)?\s*\{\s*$`
	MTHD            = "method" // match methods on a type
	RGX_MTHD        = `^func\s+\([A-Za-z_]+\s+([^)]*)\)\s*?([A-Za-z_]\w*)\s*\(([^)]*)\)\s*\(?([^{)]*|[^{\s]+)?\)?\s*\{\s*$`
	STRUCT          = "struct" // match struct types
	RGX_STRUCT      = `^type\s+(.+?)\s*struct\s*{\s*$`
	STRUCTFLD       = "structFld" // match fields of a struct type
	RGX_STRUCT_FLD  = `^\s*([A-Za-z_]\w*)\s+([*\[\]\w.{}]+)\s*([` + "`" + `].*[` + "`" + `]|\/\/.*)?\s*$`
	ENDSTMNT        = "structEnd" // match the END of a struct type - only '}'
	RGX_ENDSTMNT    = `^\s*}\s*$`
	CLOSE_PAREN     = "closed_paren"
	RGX_CLOSE_PAREN = `^\)\s*$`
	TYPEMAP         = "typeMap" // match map types (map aliased to a type)
	RGX_TYPE_MAP    = `^type\s+(\w+)\s+(map\[.*)$`
)

// map patterns to find type
type RgxPatternMap map[string]string

func MapRegexPatterns() RgxPatternMap {
	return RgxPatternMap{
		EMPTY:       RGX_EMPTY,
		PKG:         RGX_PKG,
		IMP:         RGX_IMP,
		IMPKG:       RGX_IMPKG,
		CMNT:        RGX_CMNT,
		FUNC:        RGX_FUNC,
		MTHD:        RGX_MTHD,
		STRUCT:      RGX_STRUCT,
		STRUCTFLD:   RGX_STRUCT_FLD,
		CLOSE_PAREN: RGX_CLOSE_PAREN,
		ENDSTMNT:    RGX_ENDSTMNT,
		TYPEMAP:     RGX_TYPE_MAP,
	}

}

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

// compile all regex patterns in the RegexPtrns type - map returned Regexp type
// to the item name
func CompileRgx(rm RgxPatternMap) *RgxReady {
	rr := RgxReady{}
	for name, ptrn := range rm {
		rr[name] = regexp.MustCompile(ptrn)
	}
	return &rr
}
