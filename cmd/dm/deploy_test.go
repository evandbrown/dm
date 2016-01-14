package main

import (
	"strings"
	"testing"
)

func TestGetName(t *testing.T) {
	name, err := getName(true)
	if name == "" || err != nil {
		t.Errorf("Error retrieving name: %s", err.Error())
	}
	if !strings.Contains(name, "-") {
		t.Errorf("Name should contain a uuid: %s", name)
	}
	if strings.Contains(name, "/") {
		t.Errorf("Name should not contain a /: %s", name)
	}

	uid = false
	name, err = getName(false)
	if name == "" || err != nil {
		t.Errorf("Error retrieving name: %s", err.Error())
	}
	if strings.Contains(name, "-") {
		t.Errorf("Name should not contain a uuid: %s", name)
	}
}
