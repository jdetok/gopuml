package conf

import (
	"encoding/json"
	"os"

	"github.com/jdetok/gopuml/pkg/errd"
)

// read the .gopuml.json config file

// read from JSON config file
type Conf struct {
	CnfPath     string
	ProjectName string   `json:"project_name"`
	ProjectRoot string   `json:"project_root"`
	ExcludeDirs []string `json:"exclude_dirs"`
	PumlOut     string   `json:"puml_out"`
}

func GetConf(fname string) (*Conf, error) {
	var gp Conf
	if err := gp.OpenConfigF(fname); err != nil {
		return nil, err
	}
	return &gp, nil
}

// attempt read the file at fname - if doesn't exist, call CreateConf
func (cnf *Conf) OpenConfigF(fname string) error {
	f, err := os.Open(fname)
	if err == nil {
		if err := json.NewDecoder(f).Decode(cnf); err != nil {
			return &errd.ConfDecodeError{Path: fname, Err: err}
		}
		defer f.Close()
		cnf.CnfPath = f.Name()
		return nil
	}
	return cnf.CreateConf(fname)
}

// conf file exists - decode the JSON to JSONConf struct fields
func (cnf *Conf) ReadConf(f *os.File) error {
	return json.NewDecoder(f).Decode(cnf)
}

// conf file does not exist - create one in the root dir
func (cnf *Conf) CreateConf(fname string) error {
	f, err := os.Create(fname)
	if err != nil {
		return &errd.FileCreateError{Path: fname, Err: err}
	}
	defer f.Close()

	// get bytes of encoded and indented JSON
	b, err := json.MarshalIndent(cnf, "", "    ")
	if err != nil {
		return &errd.JSONEncodeError{FName: fname, Err: err}
	}

	// write json bytes to JSON conf
	if _, err := f.Write(b); err != nil {
		return &errd.WriterError{WriterLoc: f.Name(), NumBytes: len(b), Err: err}
	}
	cnf.CnfPath = f.Name()
	return nil
}
