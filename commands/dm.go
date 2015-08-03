package commands

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/conf"
	"github.com/spf13/cobra"
)

var Verbose bool
var Project, Name string
var DmCmd = &cobra.Command{
	Use: "dm",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logging()
	},
}

func init() {
	DmCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	DmCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "Google Cloud Platform project name")
	DmCmd.PersistentFlags().StringVarP(&Name, "name", "n", "", "Name of the deployment to use")
}

func Execute() {
	addCommands()
	DmCmd.Execute()
}

func addCommands() {
	DmCmd.AddCommand(deployCmd)
	DmCmd.AddCommand(deleteCmd)
}

func logging() {
	if Verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func checkConfig() {
	config, err := conf.ReadDeploymentConfig()
	if err != nil {
		log.Fatal("No know deployments found. Use `dm deploy` to create a new one")
	}

	if len(config.Deployments) > 1 && Name == "" {
		log.Fatal("Multiple deployments found. Use the --name flag to specify the deployment to use with the command.")
	}

	if _, ok := config.Deployments[Name]; Name != "" && !ok {
		log.Fatal(fmt.Sprintf("Deployment name '%s' not found", Name))
	}
}
