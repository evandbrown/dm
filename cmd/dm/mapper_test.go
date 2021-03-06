package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestVarsDotYAMLMapper(t *testing.T) {
	cases := []struct {
		Input       string
		Result      map[string]string
		ExpectError bool
	}{
		{
			"<",
			nil,
			true,
		},
		{
			"variables:\n  one: two\n  three: four",
			map[string]string{
				"one":   "two",
				"three": "four",
			},
			false,
		},
	}

	for _, c := range cases {
		vm := VarsDotYAMLMapper{}
		err := vm.Parse(strings.NewReader(c.Input))

		if err == nil && c.ExpectError {
			t.Errorf("expected error parsing %v\n\n%s", c.Input, err)
		}

		actual := vm.Map()
		if !reflect.DeepEqual(actual, c.Result) {
			t.Fatalf("\nexpected: %v\n\ngot: %#v", actual, c.Result)
		}
	}
}
