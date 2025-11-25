package conf

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/jdetok/gopuml/cli"
	"github.com/jdetok/gopuml/pkg/errd"
)

// command line .gopuml.json setup - init flag must be passed
// create config file from user input, return pointer to Conf struct and write
// the generated .gopuml.json file to the project root
func CLISetup(fname string) (*Conf, error) {
	cnf := Conf{}

	next := cli.AskBool(
		"user passed init flag - starting .gopuml.json CLI setup\n"+
			"press enter (return) to continue or N to exit\n...", "n")

	if !next {
		fmt.Println("user exit")
		return nil, &errd.UserExit{}
	}
	fmt.Print("continuing with .gopuml.json CLI setup...\n\n")

	// prompt user to provide value for each field in conf struct
	if err := cnf.FieldsFromUser(); err != nil {
		return &cnf, err
	}

	//
	if err := cnf.CreateJson(fname); err != nil {
		return &cnf, err
	}

	return &cnf, nil
}

// iterate through the fields in the Conf struct
// for each field, if a json tag is found with any string
// other than "" or "-", the user is prompted to provide a value for the field
// the reflect package is used to set the corresponding input value as the value
// in the reflect.StructField
func (cnf *Conf) FieldsFromUser() error {
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
		if confFld == "" || confFld == "-" {
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
