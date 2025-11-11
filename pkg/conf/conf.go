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
	ProjectName string   `json:"project_name"`
	ProjectRoot string   `json:"project_root"`
	ExcludeDirs []string `json:"exclude_dirs"`
	PumlOut     string   `json:"puml_out"`
}

func GetConf(fname string) (*GoPumlConf, error) {
	var gp GoPumlConf
	f, err := gp.OpenConfigF(fname)
	if err != nil {
		return nil, err
	}
	gp.ConfFile = f
	// gp.ExcludeDirs = []string{}
	return &gp, nil
}

func (gp *GoPumlConf) OpenConfigF(fname string) (*os.File, error) {
	var f *os.File
	f, err := os.Open(fname)
	if err == nil {
		if err := gp.JSONConf.ReadConf(f); err != nil {
			return f, &errd.ConfDecodeError{Path: fname, Err: err}
		}
		return f, nil
	}
	return gp.JSONConf.CreateConf(fname)
}

func (jc *JSONConf) ReadConf(f *os.File) error {
	return json.NewDecoder(f).Decode(jc)
}

func (jc *JSONConf) CreateConf(fname string) (*os.File, error) {

	f, err := os.Create(fname)
	if err != nil {
		return nil, &errd.FileCreateError{Path: fname, Err: err}
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
