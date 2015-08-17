package template

import (
	"bytes"
	"fmt"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Context is the context that an interpolation is done in. This defines
// things such as available variables.
type Context struct {
	Data string
	Vars map[string]string
}

func Render(c *Context) (bytes.Buffer, error) {
	log.Debugf("Rendering:\n%v\n\nwith vars:\n%v\n\n", c.Data, c.Vars)
	var result bytes.Buffer
	tmpl, err := template.New("root").Funcs(Funcs(c)).Parse(c.Data)
	if err != nil {
		return result, err
	}

	err = tmpl.Execute(&result, "")
	if err != nil {
		return result, err
	}
	return result, nil
}

type ContextBuilder struct {
	// Template data to be rendered
	Data string

	// Variables is the mapping of user variables that the
	// "user" function reads from.
	UserVars        map[string]string
	AllowedUserVars map[string]string
	EnableEnv       bool
	Path            string

	Error error
}

func NewContextBuilder() *ContextBuilder {
	c := &ContextBuilder{
		UserVars:        make(map[string]string),
		AllowedUserVars: make(map[string]string),
		EnableEnv:       true,
	}

	return c
}

func (c *ContextBuilder) RenderConfig() (*Config, error) {
	log.Debugf("rendering configuration at %v", c.Path)
	config := &Config{}
	config.Path = c.Path

	ctx, err := c.context()
	if err != nil {
		return config, err
	}

	log.Debugf("got context with data:\n\n%v\n\n and vars:\n\n%v\n", ctx.Data, ctx.Vars)
	result, err := Render(ctx)
	log.Debugf("rendered context:\n\n%v", result.String())
	if err != nil {
		return config, err
	}
	config.Raw = result.Bytes()

	err = yaml.Unmarshal(result.Bytes(), config)
	if err != nil {
		return config, err
	}
	log.Debugf("marshalled config to yaml")
	return config, nil
}

func (c *ContextBuilder) context() (*Context, error) {
	err := c.Validate()
	if err != nil {
		return &Context{}, err
	}
	return &Context{
		Data: c.Data,
		Vars: c.UserVars,
	}, nil
}

func (c *ContextBuilder) AddUserVarsFromProvider(vp *VarProvider) {
	if c.Error != nil {
		return
	}

	vars, err := vp.Mapper.Map(vp.Source)
	if err != nil {
		c.Error = err
		return
	}
	log.Debugf("adding user vars from provider:\n\n%v\n", vars)
	c.AddUserVars(vars, vp.Constrain)
}

func (c *ContextBuilder) AddUserVars(vars map[string]string, constrain bool) {
	log.Debugf("adding user vars. contrained=%v", constrain)
	for k, v := range vars {
		if constrain {
			c.AllowedUserVars[k] = v
			c.AddUserVar(k, v)
		} else {
			c.AddUserVar(k, v)
		}
	}
	log.Debugf("user vars now: %v\n", c.UserVars)
}

func (c *ContextBuilder) AddUserVar(k string, v string) {
	if c.EnableEnv {
		log.Debugf("--rendering env var \"%v\"\n", v)
		rendered, err := Render(&Context{Data: v})
		if err != nil {
			c.Error = err
		}
		log.Debugf("--rendered: \"%v\"\n", rendered.String())
		v = rendered.String()
	}
	log.Debugf("--adding user var %v with value %s\n", k, v)
	c.UserVars[k] = v
}

func (c *ContextBuilder) Validate() error {
	for k, v := range c.UserVars {
		if _, ok := c.AllowedUserVars[k]; !ok {
			return fmt.Errorf("variable %s was provided but not in the allowed list", k)
		}
		if v == "" {
			return fmt.Errorf("variable %s is empty", k)
		}
	}
	log.Debugf("Validated vars: %v\n", c.UserVars)
	return nil
}
