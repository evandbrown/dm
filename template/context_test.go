package template

import (
	"os"
	"reflect"
	"testing"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetLevel(log.InfoLevel)
}

func TestUserVars(t *testing.T) {
	ctxb := NewContextBuilder()
	userVars := map[string]string{
		"foo": "bar",
	}
	ctxb.AddUserVars(userVars, false)
	if !reflect.DeepEqual(ctxb.UserVars, userVars) {
		t.Fatalf("\nexpected: %v\n\ngot: %#v", ctxb.UserVars, userVars)
	}
}

func TestUserVarsConstrain(t *testing.T) {
	ctxb := NewContextBuilder()
	userVars := map[string]string{
		"foo": "bar",
	}
	ctxb.AddUserVars(userVars, true)

	ctxb.AddUserVars(userVars, false)
	if !reflect.DeepEqual(ctxb.UserVars, userVars) {
		t.Fatalf("\nexpected: %v\n\ngot: %#v", ctxb.UserVars, userVars)
	}
}

func TestUserVarsFromEnv(t *testing.T) {
	os.Setenv("foo", "bar")
	ctxb := NewContextBuilder()
	userVars := map[string]string{
		"foo": "{{env `foo`}}",
	}
	expected := map[string]string{
		"foo": "bar",
	}
	ctxb.AddUserVars(userVars, false)
	// should not equal
	if reflect.DeepEqual(ctxb.UserVars, userVars) {
		t.Fatalf("\nexpected: %v\n\ngot: %#v", ctxb.UserVars, userVars)
	}
	if !reflect.DeepEqual(ctxb.UserVars, expected) {
		t.Fatalf("\nexpected: %v\n\ngot: %#v", ctxb.UserVars, userVars)
	}
}