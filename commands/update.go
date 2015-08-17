package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/template"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing configuration in Deployment Manager.",
}

func init() {
	updateCmd.PreRun = func(cmd *cobra.Command, args []string) {
		requireName()
	}
	updateCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(update(cmd, args))
	}
}

func update(cmd *cobra.Command, args []string) error {
	// Get config from disk
	config, err := conf.ReadDeploymentConfig()
	if err != nil {
		return err
	}
	c := config.Deployments[Name]

	deployment, err := googlecloud.GetDeployment(c.Project, c.Id)
	if err != nil {
		return err
	}

	log.Infof("Updating deployment %s", Name)

	d, err := template.GenerateDeployment(Name, "", &template.Config{})
	if err != nil {
		return err
	}
	d.Intent = "UPDATE"
	d.Fingerprint = deployment.Fingerprint

	service, err := googlecloud.GetService()
	if err != nil {
		return err
	}
	_, err = service.Deployments.Update(Project, Name, d).Do()
	if err != nil {
		return err
	}
	fmt.Printf("Updated deployment %s.\n", Name)
	return nil
}