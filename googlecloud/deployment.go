package googlecloud

import (
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/util"
	"github.com/google/google-api-go-client/deploymentmanager/v2beta2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/yaml.v2"
)

func GetService() (*deploymentmanager.Service, error) {
	client, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform")
	service, err := deploymentmanager.New(client)
	util.Check(err)

	return service, nil
}

func NewDeployment(name string, description string, configPath string) *deploymentmanager.Deployment {

	tc, err := ParseConfig(configPath)
	util.Check(err)

	d := &deploymentmanager.Deployment{Name: name, Description: description, Target: tc}
	return d
}

func ParseConfig(configPath string) (*deploymentmanager.TargetConfiguration, error) {
	log.Debugf("Parsing deployment configuration from %s", configPath)
	cfg, err := ioutil.ReadFile(configPath)
	util.Check(err)

	c := C{}
	err = yaml.Unmarshal(cfg, &c)
	util.Check(err)
	log.Debugf("Parsing imports %v", c)
	imports := make([]*deploymentmanager.ImportFile, len(c.Imports))

	for i, imp := range c.Imports {
		log.Debugf("Parsing import %v", imp)
		importPath := strings.Split(configPath, "/")
		importPath[len(importPath)-1] = imp.Path
		f := strings.Join(importPath, "/")
		log.Debugf("Reading imported file %s", f)
		impBytes, err := ioutil.ReadFile(f)
		util.Check(err)
		imports[i] = &deploymentmanager.ImportFile{Name: imp.Path, Content: string(impBytes)}
	}

	tc := &deploymentmanager.TargetConfiguration{Config: string(cfg), Imports: imports}

	return tc, nil
}

type C struct {
	Imports []struct{ Path string }
}
