package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/util"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List deployments",
}

func init() {
	lsCmd.PreRun = func(cmd *cobra.Command, args []string) {
		requireConfig()
	}
	lsCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(ls(cmd, args))
	}
}

func ls(cmd *cobra.Command, args []string) error {
	// Get config from disk
	config, _ := conf.ReadDeploymentConfig()

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Deployment Name\tProject\tState\t")
	for _, c := range config.Deployments {
		_, err := googlecloud.GetDeployment(c.Project, c.Id)
		if err != nil {
			return err
		}
		//TODO: show deployment status
		fmt.Fprintf(w, "%s\t%s\t%s\t\n", c.Id, c.Project, "")
	}
	w.Flush()
	return nil
}
