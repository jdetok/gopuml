package cli

import (
	"flag"
)

const (
	// flag to set location/name of config json file
	CONF         = "conf"
	DEFAULT_CONF = ".gopuml.json"

	// flag to set project root location
	ROOT         = "root"
	DEFAULT_ROOT = "."

	// flag to set file type (currently only works with Go)
	FTYP         = "ftyp"
	DEFAULT_FTYP = ".go"
)

type ArgMap map[string]*string

// parse flag args
func MapArgs() *ArgMap {
	root := ROOT
	conf := CONF
	ftyp := FTYP

	var args = ArgMap{
		ROOT: &root,
		CONF: &conf,
		FTYP: &ftyp,
	}

	flag.StringVar(args[ROOT], ROOT, DEFAULT_ROOT,
		"root dir of project gopuml is called on")

	flag.StringVar(args[CONF], CONF, DEFAULT_CONF,
		"conf file location (default project root)")

	flag.StringVar(args[FTYP], FTYP, DEFAULT_FTYP,
		"file type to look for")

	flag.Parse()

	return &args
}
