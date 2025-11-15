package rgx

import (
	"regexp"
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

type RgxPatternMap map[string]string

// think for each file i need to do a struct with arrays of funcs, structs, etc
type Rgx struct {
	ptrns   RgxPatternMap
	ready   RgxReady
	PkgMap  RgxPkgMap
	RDirMap RgxDirMap
	Funcs   []*RgxFunc
	Structs []*RgxStruct
}

func NewRgx() *Rgx {
	return &Rgx{
		ptrns:   RgxPatternMap{},
		ready:   *CompileRgx(),
		PkgMap:  RgxPkgMap{},
		RDirMap: RgxDirMap{},
	}
}

type RgxFile struct {
	Pkg     string // package (dir) it belongs to
	RgxPkg  string
	Name    string // file
	Structs []*RgxStruct
	Funcs   []*RgxFunc
	Methods []*RgxFunc
}

// map regexp objext to item name (can be used through runtime to match to)
type RgxReady map[string]*regexp.Regexp

// map package string to file map
type RgxDirMap map[string]RgxFileMap
type RgxPkgMap map[string]RgxFileMap

// map file name to rgx match type, each filemap is mapped to a package name
type RgxFileMap map[string]*RgxFile

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
	Fields []RgxStructFld
}

type RgxStructFld struct {
	Name  string // field name
	DType string // field type
	Cmnt  string
	Tag   string
}
