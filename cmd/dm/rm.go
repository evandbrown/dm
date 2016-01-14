package main

import (
	"fmt"

	"github.com/evandbrown/dm/googlecloud"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a deployment",
}

func init() {
	deleteCmd.PreRun = func(cmd *cobra.Command, args []string) {
		requireName()
	}

	deleteCmd.Run = func(cmd *cobra.Command, args []string) {
		Check(rm(cmd, args))
	}
}

func rm(cmd *cobra.Command, args []string) error {
	service, err := googlecloud.GetService()
	Check(err)

	config, err := ReadDeploymentConfig()
	Check(err)

	d := config.Deployments[Name]
	call := service.Deployments.Delete(d.Project, d.Id)
	_, err = call.Do()
	Check(err)

	err = RemoveDeployment(Name)
	Check(err)

	fmt.Printf("Deleted deployment %s\n", Name)
	return nil
}
