package conf

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jdetok/gopuml/cli"
)

// create new gopuml file if one doesn't exist
func Setup(fname string) *Conf {
	cnf := Conf{}

	next := cli.AskBool(`gopuml called with --init flag
press enter (return) to configure .gopuml.json
press N to exit
	 `, "n")

	if !next {
		fmt.Println("user exit")
		return nil
	}
	fmt.Println("continue with CLI setup implementation")
	cnf.CLIFieldSetup()
	cnf.CreateConf(fname)
	// if not N, ask whether to create template file or continue with CLI setup
	// setup - ask for each value
	// TODO: better CLI function for console prommpt

	return &cnf
}

func (cnf *Conf) CLIFieldSetup() error {
	// check that value is addressable
	cv := reflect.ValueOf(cnf)
	if cv.Kind() != reflect.Pointer || cv.IsNil() {
		return fmt.Errorf("FillFromConsole requires a non-nil *Conf")
	}

	// get value and type of element cnf is pointing to
	cv = cv.Elem()
	ct := cv.Type()

	// loop through fields in conf struct
	for i := range cv.NumField() {
		fld := ct.Field(i) // struct field
		val := cv.Field(i) // value in struct field

		// skip if not settable
		if !val.CanSet() {
			continue
		}

		// get the name of the json field the structfield encodes to
		confFld := fld.Tag.Get("json")
		if confFld == "" {
			continue
		}

		// ask user for value to fill field
		cliVal := cli.AskStr(fmt.Sprintf(
			"enter a value or enter (return) to set default %s...\n%s: ",
			confFld, confFld))

		// assign the input value to the struct field given its kind
		switch val.Kind() {
		case reflect.String:
			val.SetString(strings.TrimSpace(cliVal))
			fmt.Printf("successfully set %s as type string to field %s\n\n",
				cliVal, confFld)
		case reflect.Slice:
			// accept comma separated string for slices, split to slice of strs
			subs := strings.Split(cliVal, ",")
			for s := range subs {
				// delete index if it's blank
				// (handles input like ".go, .py," with trailing comma)
				if subs[s] == "" {
					subs = slices.Delete(subs, s, s+1)
					continue
				}
				// otherwise, remove any trailing space and replace string at idx
				subs[s] = strings.TrimSpace(subs[s])
			}
			// set val of struct field as value of the slice
			toSet := reflect.ValueOf(subs)
			val.Set(toSet)
			fmt.Printf("successfully set %s as type %s to field %s\n\n",
				cliVal, toSet.Kind(), confFld)
		default:
			continue
		}
	}
	return nil
}
