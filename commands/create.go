package commands

import (
	"fmt"

	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var uuid bool
var manifest, name string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment",
}

func init() {
	createCmd.Flags().StringVarP(&manifest, "manifest", "m", "manifest.yaml", "The name of the manifest to deploy.")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the deployment. Defaults to the name of the currect directory")
	createCmd.Flags().BoolVarP(&uuid, "uuid", "u", true, "Should a 7 char UUID be appended to deployment name. Defaults is yes")
	createCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(create(cmd, args))
	}
}

func create(cmd *cobra.Command, args []string) error {
	fmt.Println("Creating...")

	f := cmd.Flags().Lookup("verbose")
	fmt.Println(f.Value)

	service, err := util.GetService()
	if err != nil {
		return err
	}

	util.NewDeployment("foo", "bar", "/tmp/config.yaml")
	service.Deployments.List("majestic-device-95716")
	return nil
}
