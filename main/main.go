package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/conf"
	"github.com/jdetok/gopuml/pkg/dir"
	"github.com/jdetok/gopuml/pkg/errd"
	"github.com/jdetok/gopuml/pkg/puml"
	"github.com/jdetok/gopuml/pkg/rgx"
)

func main() {
	var err error
	var userExit *errd.UserExit // used for graceful shutdown
	var cnf *conf.Conf

	// get args passed with program execution
	args := *cli.ParseArgs()

	// root flag
	rootDir := *args.ArgMap[cli.ROOT]

	// ftyp flag
	ftyp := *args.ArgMap[cli.FTYP]

	// .gopuml.conf file location
	confFile := *args.ArgMap[cli.CONF]
	confPath := filepath.Join(rootDir, confFile)

	// if init flag is passed, run cli setup to create file
	if args.Init {
		// setup file
		cnf, err = conf.CLISetup(confPath)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("conf file successfully created at", cnf.CnfPath)
	} else {
		// read or .gopuml.json conf file
		cnf, err = conf.DecodeJsonConf(confPath)
		if err != nil {
			log.Fatal(err)
		}
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
	fmt.Printf("successfully parsed regular expressions in %s files\n", cnf.Langs[0])

	p := puml.Puml{
		Dgm: puml.NewUmlClass(fmt.Sprintf("%s | %s", cnf.ProjectName, cnf.ClassT), r),
	}
	fmt.Println(filepath.Join(cnf.ProjectRoot, cnf.OutDir))

	err = p.WriteOutput(filepath.Join(cnf.ProjectRoot, cnf.OutDir), cnf.ClassF)
	if err != nil {
		if errors.As(err, &userExit) {
			fmt.Println(err)
			os.Exit(0)
		}
		log.Fatal(err)
	}
}
