package main

import (
	"fmt"
	"log"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/conf"
	"github.com/jdetok/gopuml/pkg/dir"
)

func main() {
	args := cli.ParseArgs()

	fName := args.Root[1] + "/" + args.ConfFile[1]
	fType := args.FType[1]

	_, err := conf.NewGoPumlConf(fName)
	if err != nil {
		log.Fatal(err)
	}
	files := map[string]string{}
	count, err := dir.CheckDirForFType(args.Root[1], fType, files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d total %s files\n", count, fType)
	fmt.Println(files)
}
