package commands

import (
	"fmt"

	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var statCmd = &cobra.Command{
	Use:   "stat",
	Short: "Stat a deployment, listing its resources",
}

func init() {
	statCmd.PreRun = func(cmd *cobra.Command, args []string) {
		checkConfig()
	}
	statCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(stat(cmd, args))
	}
}

func stat(cmd *cobra.Command, args []string) error {
	service, err := googlecloud.GetService()
	util.Check(err)

	call := service.Resources.List(Project, Name)
	resources, error := call.Do()
	util.Check(error)
	for _, r := range resources.Resources {
		fmt.Printf("%s\t%s\n", r.Type, r.Name)
	}
	return nil
}
