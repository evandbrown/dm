package commands

import (
	"fmt"

	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List deployments",
}

func init() {
	lsCmd.PreRun = func(cmd *cobra.Command, args []string) {
		checkConfig()
	}
	lsCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(ls(cmd, args))
	}
}

func ls(cmd *cobra.Command, args []string) error {
	config, _ := conf.ReadDeploymentConfig()

	for _, c := range config.Deployments {
		fmt.Printf("%s\n", c.Id)
	}
	return nil
}
