package util

import (
	"bytes"
	"github.com/google/google-api-go-client/deploymentmanager/v2beta2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"io"
	"os"
)

func GetService() (*deploymentmanager.Service, error) {
	client, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform")
	service, err := deploymentmanager.New(client)
	Check(err)

	return service, nil
}

func NewDeployment(name string, description string, confAbsPath string) *deploymentmanager.Deployment {
	f, err := os.Open(confAbsPath)
	Check(err)

	tc, err := ParseConfig(f)
	Check(err)

	d := &deploymentmanager.Deployment{Name: name, Description: description, Target: tc}
	return d
}

func ParseConfig(config io.Reader) (*deploymentmanager.TargetConfiguration, error) {
	buf := new(bytes.Buffer)
	buf.ReadFrom(config)
	s := buf.String()

	tc := &deploymentmanager.TargetConfiguration{Config: s}

	return tc, nil
}
