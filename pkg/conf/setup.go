package conf

import (
	"fmt"

	"github.com/jdetok/gopuml/cli"
)

// create new gopuml file if one doesn't exist
func Setup() *Conf {
	cnf := Conf{}

	next := cli.Ask(`gopuml called with --init flag
press enter (return) to configure .gopuml.json
press N to exit
	 `, "n")

	if !next {
		fmt.Println("user exit")
		return nil
	}
	fmt.Println("continue with CLI setup implementation")
	// if not N, ask whether to create template file or continue with CLI setup
	// setup - ask for each value
	// TODO: better CLI function for console prommpt

	return &cnf
}
