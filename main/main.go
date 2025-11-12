package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/conf"
	"github.com/jdetok/gopuml/pkg/dir"
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

	// fmt.Println(cnf)
	// get the passed file type, recursively loop through root to find files
	// with that type
	dirMap, err := dir.MapFiles(rootDir, ftyp, cnf.ExcludeDirs)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("%d dirs within \"%s\" contain %s files\n", len(*dirMap),
	// 	rootDir, ftyp)
	// fmt.Println(*dirMap)
	f, err := dirMap.OpenFile("pkg/dir", "mapf.go")
	if err != nil {
		log.Fatal(err)
	}
	rgxPtrns := rgx.NewRgxPtrns()
	rgxReady, err := rgxPtrns.CompileRgx()
	if err != nil {
		log.Fatal(err)
	}
	if err := rgxReady.RgxParseFile(f); err != nil {
		log.Fatal(err)
	}
}
