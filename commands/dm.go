package commands

import (
	"github.com/spf13/cobra"
)

var Verbose bool
var Project string
var DmCmd = &cobra.Command{
	Use: "dm",
}

func init() {
	DmCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	DmCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "Google Cloud Platform project name")
}

func Execute() {
	AddCommands()
	DmCmd.Execute()
}

func AddCommands() {
	DmCmd.AddCommand(createCmd)
}
