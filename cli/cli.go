package cli

import (
	"flag"
)

const DEFAULT_CONF_FILE = ".gopuml.json"
const FTYPE_GO = ".go"

type Args struct {
	Root     [2]string // root of the project where gopuml is run
	ConfFile [2]string // file location of .gopuml.json (default project root)
	FType    [2]string // file type (currently only works with .go)
}

// parse flag args
func ParseArgs() *Args {
	var args = Args{
		Root:     [2]string{"root", ""},
		ConfFile: [2]string{"conf", ""},
		FType:    [2]string{"ftype", ""},
	}

	// default to config file in the root directory
	flag.StringVar(&args.ConfFile[1], "conf", DEFAULT_CONF_FILE,
		"conf file location (default project root)")

	flag.StringVar(&args.Root[1], "root", "",
		"root dir of project gopuml is called on")

	flag.StringVar(&args.FType[1], "ftype", FTYPE_GO, "file type to look for")

	flag.Parse()

	return &args
}
