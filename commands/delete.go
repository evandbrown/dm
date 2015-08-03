package commands

import (
	log "github.com/Sirupsen/logrus"
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

	d := config.Deployments[Name]
	call := service.Deployments.Delete(d.Project, d.Id)
	_, err = call.Do()
	util.Check(err)

	err = conf.RemoveDeployment(Name)
	util.Check(err)

	log.Printf("Deleted deployment %s", Name)
	return nil
}
