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

	// bool var for init
	INIT         = "init"
	DEFAULT_INIT = false
)

type ArgMap map[string]*string

type Args struct {
	ArgMap
	Init bool
}

func ParseArgs() *Args {
	var args Args

	// get string flags
	args.ArgMap = *MapArgs()

	// setup bool init flag
	var initBool bool
	flag.BoolVar(&initBool, INIT, DEFAULT_INIT,
		"gopuml init - creates conf file")

	// parse flags
	flag.Parse()

	// set bool init val
	args.Init = initBool
	return &args
}

// parse flag args
func MapArgs() *ArgMap {
	root := ROOT
	conf := CONF
	ftyp := FTYP
	init := INIT

	// setup flag names
	var argMap = ArgMap{
		ROOT: &root,
		CONF: &conf,
		FTYP: &ftyp,
		INIT: &init,
	}

	flag.StringVar(argMap[ROOT], ROOT, DEFAULT_ROOT,
		"root dir of project gopuml is called on")

	flag.StringVar(argMap[CONF], CONF, DEFAULT_CONF,
		"conf file location (default project root)")

	flag.StringVar(argMap[FTYP], FTYP, DEFAULT_FTYP,
		"file type to look for")

	return &argMap
}
