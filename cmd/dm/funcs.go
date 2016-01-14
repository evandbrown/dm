package main

import (
	"os"
	"strconv"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Funcs are the interpolation funcs that are available within interpolations.
var FuncGens = map[string]FuncGenerator{
	"env":       funcGenEnv,
	"timestamp": funcGenTimestamp,
	"var":       funcGenVar,
}

// InitTime is the UTC time when this package was initialized. It is
// used as the timestamp for all configuration templates so that they
// match for a single build.
var InitTime time.Time

func init() {
	InitTime = time.Now().UTC()
}

// FuncGenerator is a function that given a context generates a template
// function for the template.
type FuncGenerator func(*Context) interface{}

// Funcs returns the functions that can be used for interpolation given
// a context.
func Funcs(ctx *Context) template.FuncMap {
	result := make(map[string]interface{})
	for k, v := range FuncGens {
		log.Debugf("adding %v handler func", k)
		result[k] = v(ctx)
	}

	return template.FuncMap(result)
}

func funcGenEnv(ctx *Context) interface{} {
	return func(k string) (string, error) {
		log.Debugf("rendering env var for '%v'", k)
		return os.Getenv(k), nil
	}
}

func funcGenTimestamp(ctx *Context) interface{} {
	return func() string {
		return strconv.FormatInt(InitTime.Unix(), 10)
	}
}

func funcGenVar(ctx *Context) interface{} {
	return func(k string) string {
		if ctx == nil || ctx.Vars == nil {
			return ""
		}

		return ctx.Vars[k]
	}
}
