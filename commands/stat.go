package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

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
		requireName()
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)
	fmt.Fprintln(w, "Resource Type\tName\t")
	for _, r := range resources.Resources {
		fmt.Fprintf(w, "%s\t%s\t\n", r.Type, r.Name)
	}
	w.Flush()
	return nil
}
