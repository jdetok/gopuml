package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/conf"
	"github.com/jdetok/gopuml/pkg/dir"
	"github.com/jdetok/gopuml/pkg/puml"
	"github.com/jdetok/gopuml/pkg/rgx"
)

func main() {
	// get args passed with program execution
	args := *cli.ParseArgs()
	rootDir := *args.ArgMap[cli.ROOT]
	confFile := *args.ArgMap[cli.CONF]
	ftyp := *args.ArgMap[cli.FTYP]

	fmt.Println("init flag:", args.Init)

	// read the args for root and conf, join as filepath,
	// read or create .gopuml.json conf file
	confPath := filepath.Join(rootDir, confFile)
	cnf, err := conf.GetConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	if args.Init {
		fmt.Println("conf file successfully created at", cnf.CnfPath)
	} else {
		fmt.Printf("conf file at %s successfully read and decoded into conf struct\n",
			cnf.CnfPath)
	}

	// get the passed file type, recursively loop through root to find files
	// with that type
	dirMap, err := dir.MapFiles(rootDir, ftyp, cnf.ExcludeDirs)
	if err != nil {
		log.Fatal(err)
	}

	// compile regex patterns
	r := rgx.NewRgx()

	// go through each file in the dirmap and match to regex patterns
	if err := r.Parse(dirMap); err != nil {
		log.Fatal(err)
	}
	for _, fn := range r.Funcs {
		fmt.Println(*fn)
	}
	for _, st := range r.Structs {
		fmt.Println(st)
	}

	p := puml.Puml{
		Dgm: &puml.UmlClass{Title: "Test Class Diagram"},
	}
	if err := p.WriteOutput(cnf.OutDir, cnf.ClassF); err != nil {
		log.Fatal(err)
	}
}
