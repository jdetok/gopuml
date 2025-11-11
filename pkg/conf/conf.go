package conf

import (
	"encoding/json"
	"fmt"
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
	gp.ExcludeDirs = []string{}
	return &gp, nil
}

func (gp *GoPumlConf) ReadConf(f *os.File) {
	buf := []byte{200}
	f.Read(buf)
	fmt.Print(string(buf))
}

func (gp *GoPumlConf) OpenConfigF(fname string) (*os.File, error) {
	var f *os.File
	f, err := os.Open(fname)
	if err == nil {
		// gp.ReadConf(f)
		return f, nil
	}
	return gp.JSONConf.CreateConfFile(fname)
}

func (jc *JSONConf) CreateConfFile(fname string) (*os.File, error) {

	f, err := os.Create(fname)
	if err != nil {
		return nil, &errd.FileCreateError{FName: fname, Err: err}
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
