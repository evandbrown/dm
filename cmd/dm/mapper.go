package main

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Implement the VarMapper interface to convert
// variable for a template into a map
type VarMapper interface {
	Parse(r io.Reader)
	Map() (map[string]string, error)
	Keys() map[string]struct{}
}

type VarsDotYAMLMapper struct {
	data map[string]string
	keys []string
}

func (v *VarsDotYAMLMapper) Parse(r io.Reader) error {
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	vars := &Variables{}
	err = yaml.Unmarshal(contents, vars)
	if err != nil {
		return err
	}
	v.data = vars.Variables

	v.keys = make([]string, len(v.data), len(v.data))
	i := 0
	for k, _ := range v.data {
		v.keys[i] = k
		i++
	}
	return nil
}

func (v *VarsDotYAMLMapper) Map() map[string]string {
	return v.data
}

func (v *VarsDotYAMLMapper) Keys() []string {
	return v.keys
}
