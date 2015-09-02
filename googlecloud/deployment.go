package googlecloud

import (
	log "github.com/Sirupsen/logrus"
	"github.com/google/google-api-go-client/deploymentmanager/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GetService() (*deploymentmanager.Service, error) {
	client, err := google.DefaultClient(oauth2.NoContext, "https://www.googleapis.com/auth/cloud-platform")
	if err != nil {
		return &deploymentmanager.Service{}, err
	}
	return deploymentmanager.New(client)
}

func GetDeployment(project string, name string) (*deploymentmanager.Deployment, error) {
	log.Debugf("Getting deployment %s from project %s", name, project)
	// Get service
	service, err := GetService()
	if err != nil {
		return nil, err
	}
	d, err := service.Deployments.Get(project, name).Do()
	if err != nil {
		return nil, err
	}
	return d, nil
}
