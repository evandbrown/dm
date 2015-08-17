package template

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/google/google-api-go-client/deploymentmanager/v2beta2"
)

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

func GenerateDeployment(name string, description string, config *Config) (*deploymentmanager.Deployment, error) {
	log.Debugf("Parsing deployment configuration from %s", config.Path)
	d := &deploymentmanager.Deployment{Name: name, Description: description}

	dmConfig := APIConfig{
		Base: config,
	}
	err := dmConfig.ParseImports()
	if err != nil {
		return d, err
	}
	tc := &deploymentmanager.TargetConfiguration{Config: string(config.Raw), Imports: dmConfig.ImportFiles}

	d.Target = tc
	return d, nil
}
