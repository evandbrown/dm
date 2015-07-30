package googlecloud

import (
	"io/ioutil"
	"strings"

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
	cfg, err := ioutil.ReadFile(configPath)
	util.Check(err)

	c := C{}
	err = yaml.Unmarshal(cfg, &c)
	util.Check(err)

	imports := make([]*deploymentmanager.ImportFile, len(c.Imports))

	for i, imp := range c.Imports {
		dirs := strings.Split(configPath, "/")
		if len(dirs) > 1 {
			dirs[len(dirs)-1] = imp.Path
		}
		impBytes, err := ioutil.ReadFile(strings.Join(dirs, "/"))
		util.Check(err)
		imports[i] = &deploymentmanager.ImportFile{Name: imp.Path, Content: string(impBytes)}
	}

	tc := &deploymentmanager.TargetConfiguration{Config: string(cfg), Imports: imports}

	return tc, nil
}

type C struct {
	Imports []struct{ Path string }
}
