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
	updateCmd.Flags().VarP(&vars, "var", "v", "A variable value to provide to the vars.yaml file for use in a deployment. Define multiple with -v var1=foo -v var2=2")
	updateCmd.PreRun = func(cmd *cobra.Command, args []string) {
		requireName()
	}
	updateCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(update(cmd, args))
	}
}

func update(cmd *cobra.Command, args []string) error {
	// Get config from disk
	dmconf, err := conf.ReadDeploymentConfig()
	if err != nil {
		return err
	}
	c := dmconf.Deployments[Name]

	depBuilder := &template.DeploymentBuilder{
		DeploymentName:  Name,
		DeploymentDesc:  "",
		ConfigFilePath:  c.Config,
		VarsDotYamlPath: c.Vars,
		CLIVars:         vars.vars,
	}

	d, err := depBuilder.GetDeployment()
	if err != nil {
		log.Warning(err)
		return err
	}

	//d.Intent = "UPDATE"
	existing, err := googlecloud.GetDeployment(c.Project, c.Id)
	if err != nil {
		return err
	}
	d.Fingerprint = existing.Fingerprint

	service, err := googlecloud.GetService()
	if err != nil {
		return err
	}
	log.Infof("Updating deployment %s", Name)
	_, err = service.Deployments.Update(Project, Name, d).Do()
	if err != nil {
		return err
	}
	fmt.Printf("Updated deployment %s.\n", Name)
	return nil
}
