package conf

import (
	"encoding/json"
	"os"

	"github.com/jdetok/gopuml/pkg/errd"
)

// read the .gopuml.json config file

type GoPumlConf struct {
	ConfFile *os.File
	JSONConf
}

// read from JSON config file
type JSONConf struct {
	ProjectName string `json:"project_name"`
	ProjectRoot string `json:"project_root"`
	PumlOut     string `json:"puml_out"`
}

// return new GoPumlConf type
func NewGoPumlConf(fname string) (*GoPumlConf, error) {
	var gp GoPumlConf
	f, err := gp.OpenConfigF(fname)
	if err != nil {
		return nil, err
	}
	gp.ConfFile = f
	return &gp, nil
}

// attempt to open config file at fname, create a new one if it doesn't exist
func (gp *GoPumlConf) OpenConfigF(fname string) (*os.File, error) {
	var f *os.File
	f, err := os.Open(fname)
	if err == nil {
		return f, nil
	}
	// create file (by marshalling JSONConf struct) if it doesn't exist
	return gp.JSONConf.CreateConfFile(fname)
}

// create new file at fName, Marshall it to json (indented for ease of use),
// write bytes to file
func (jc *JSONConf) CreateConfFile(fname string) (*os.File, error) {

	f, err := os.Create(fname)
	if err != nil {
		return nil, &errd.CreateFileError{FName: fname, Err: err}
	}
	b, err := json.MarshalIndent(jc, "", "    ")
	if err != nil {
		return nil, &errd.JSONEncodeError{FName: fname, Err: err}
	}
	if _, err := f.Write(b); err != nil {
		return nil, &errd.WriterError{WriterLoc: f.Name(), NumBytes: len(b), Err: err}
	}
	return f, nil
}
