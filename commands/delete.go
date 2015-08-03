package commands

import (
	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a deployment",
}

func init() {
	deleteCmd.PreRun = func(cmd *cobra.Command, args []string) {
		checkConfig()
	}

	deleteCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(delete(cmd, args))
	}
}

func delete(cmd *cobra.Command, args []string) error {
	service, err := googlecloud.GetService()
	util.Check(err)

	config, err := conf.ReadDeploymentConfig()
	util.Check(err)

	for _, d := range config.Deployments {
		call := service.Deployments.Delete(d.Project, d.Id)
		_, error := call.Do()
		util.Check(error)

		util.Check(conf.RemoveDeployment(Name))
		return nil
	}

	return nil
}
