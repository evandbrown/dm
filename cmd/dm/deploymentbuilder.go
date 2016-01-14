package main

import (
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/google/google-api-go-client/deploymentmanager/v2"
)

type DeploymentBuilder struct {
	DeploymentName  string
	DeploymentDesc  string
	ConfigFilePath  string
	VarsDotYamlPath string
	CLIVars         map[string]string
}

func (t *DeploymentBuilder) GetDeployment() (*deploymentmanager.Deployment, error) {
	d := &deploymentmanager.Deployment{Name: t.DeploymentName, Description: t.DeploymentDesc}

	// Create a context builder
	ctxb := NewContextBuilder()

	// Add base template
	base, err := ioutil.ReadFile(t.ConfigFilePath)
	if err != nil {
		return d, err
	}
	ctxb.Data = string(base)
	ctxb.Path = t.ConfigFilePath

	// Create a var provider for vars.yaml
	varfile, err := os.Open(t.VarsDotYamlPath)
	if err == nil {
		vp := &VarsDotYAMLMapper{}
		err = vp.Parse(varfile)
		if err != nil {
			return d, err
		}

		ctxb.AddConstraints(vp.Keys())
		ctxb.AddUserVars(vp.Map())

		// Add CLI avrs
		ctxb.AddUserVars(t.CLIVars)

	}
	// Validate and render
	err = ctxb.Validate()
	if err != nil {
		return d, err
	}

	// Check for errors
	if ctxb.Error != nil {
		return d, ctxb.Error
	}

	// Create a deployment object for the DM API
	config, err := ctxb.RenderConfig()
	if err != nil {
		return d, err
	}

	dmConfig := APIConfig{
		Base: config,
	}
	err = dmConfig.ParseImports()
	if err != nil {
		return d, err
	}
	tc := &deploymentmanager.TargetConfiguration{Config: &deploymentmanager.ConfigFile{Content: string(config.Raw)}, Imports: dmConfig.ImportFiles}

	d.Target = tc

	return d, nil
}

type APIConfig struct {
	Base *Config

	ImportFiles []*deploymentmanager.ImportFile
}

func (c *APIConfig) ParseImports() error {
	c.ImportFiles = make([]*deploymentmanager.ImportFile, len(c.Base.Imports))

	for i, imp := range c.Base.Imports {
		log.Debugf("Parsing import %v", imp)
		importPath := strings.Split(c.Base.Path, "/")
		importPath[len(importPath)-1] = imp.Path
		f := strings.Join(importPath, "/")
		log.Debugf("Reading imported file %s", f)
		impBytes, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}
		c.ImportFiles[i] = &deploymentmanager.ImportFile{Name: imp.Path, Content: string(impBytes)}
	}
	return nil
}
