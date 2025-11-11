package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/conf"
	"github.com/jdetok/gopuml/pkg/dir"
)

func main() {
	// get args passed with program execution
	args := *cli.ParseArgs()
	rootDir := *args.ArgMap["root"]
	confFile := *args.ArgMap["conf"]
	ftyp := *args.ArgMap["ftyp"]

	fmt.Println("init flag:", args.Init)

	// read the args for root and conf, join as filepath,
	// read or create .gopuml.json conf file
	confPath := filepath.Join(rootDir, confFile)
	cnf, err := conf.GetConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cnf.JSONConf)
	// get the passed file type, recursively loop through root to find files
	// with that type
	dirMap, err := dir.MapFiles(rootDir, ftyp, cnf.JSONConf.ExcludeDirs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d dirs within \"%s\" contain %s files\n", len(*dirMap),
		rootDir, ftyp)
	fmt.Println(*dirMap)

	// dirMap.ReadAll()
}
