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
	args := *cli.MapArgs()
	rootDir := *args["root"]
	confFile := *args["conf"]
	ftyp := *args["ftyp"]

	// read the args for root and conf, join as filepath,
	// read or create .gopuml.json conf file
	confPath := filepath.Join(rootDir, confFile)
	_, err := conf.GetConf(confPath)
	if err != nil {
		log.Fatal(err)
	}

	// get the passed file type, recursively loop through root to find files
	// with that type
	dirMap, err := dir.MapFiles(rootDir, ftyp)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d dirs within \"%s\" contain %s files\n", len(*dirMap),
		rootDir, ftyp)
	fmt.Println(*dirMap)
}
