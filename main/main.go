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
	// get args passed with program execution
	args := *cli.ParseArgs()
	rootDir := *args.ArgMap[cli.ROOT]
	ftyp := *args.ArgMap[cli.FTYP]
	confFile := *args.ArgMap[cli.CONF]
	confPath := filepath.Join(rootDir, confFile)

	// fmt.Println("init flag:", args.Init)
	// RUN INIT CONFIG SETUP
	if args.Init {
		// setup file
		conf.Setup(confPath)
		os.Exit(0)
	}

	// read the args for root and conf, join as filepath,
	// read or create .gopuml.json conf file
	cnf, err := conf.GetConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println("root from conf:", cnf.ProjectRoot)

	// fmt.Println("root conf:", cnf.ProjectRoot)
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
	fmt.Printf("successfully parsed regular expressions in %s files\n", cnf.Langs[0])

	p := puml.Puml{
		Dgm: puml.NewUmlClass(fmt.Sprintf("%s | %s", cnf.ProjectName, cnf.ClassT), r),
	}
	fmt.Println(filepath.Join(cnf.ProjectRoot, cnf.OutDir))

	var ue *errd.UserExit
	err = p.WriteOutput(filepath.Join(cnf.ProjectRoot, cnf.OutDir), cnf.ClassF)
	if err != nil {
		if errors.As(err, &ue) {
			fmt.Println(err)
			os.Exit(0)
		}
		log.Fatal(err)
	}
}
