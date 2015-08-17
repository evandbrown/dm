package template

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Implement the VarMapper interface to convert
// variable for a template into a map
type VarMapper interface {
	Map(r io.Reader) (map[string]string, error)
}

type VarProvider struct {
	Mapper    VarMapper
	Source    io.Reader
	Constrain bool
}

type VarsDotYAMLMapper struct{}

func (v VarsDotYAMLMapper) Map(r io.Reader) (map[string]string, error) {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	vars := &Variables{}
	err = yaml.Unmarshal(contents, vars)
	if err != nil {
		return nil, err
	}
	return vars.Variables, nil
}
