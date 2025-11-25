package conf

import (
	"fmt"

	"github.com/jdetok/gopuml/cli"
)

// create new gopuml file if one doesn't exist
func Setup() *Conf {
	cnf := Conf{}

	if cli.Ask(`gopuml called with --init flag
press enter (return) to configure .gopuml.json
press N to exit
	 `, "n") {
		fmt.Println("you said yes")
	} else {
		fmt.Println("you said no")
	}
	// if not N, ask whether to create template file or continue with CLI setup
	// setup - ask for each value
	// TODO: better CLI function for console prommpt

	return &cnf
}
