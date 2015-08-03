package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List deployments",
}

func init() {
	lsCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(ls(cmd, args))
	}
}

func ls(cmd *cobra.Command, args []string) error {
	config, err := conf.ReadDeploymentConfig()
	if err != nil {
		log.Fatal("No deployments found")
	}

	for _, c := range config.Deployments {
		fmt.Printf("%s\n", c.Id)
	}
	return nil
}
