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
	DmCmd.AddCommand(lsCmd)
	DmCmd.AddCommand(statCmd)
}

func logging() {
	log.SetFormatter(&log.TextFormatter{DisableColors: false, DisableTimestamp: true, DisableSorting: true})
	if Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

func checkConfig() {
	config, err := conf.ReadDeploymentConfig()
	if err != nil {
		log.Fatal("No deployments found. Use `dm deploy` to create a new one")
	}

	// Deployment name was provided
	if Name != "" {
		// Name not found in config file
		if _, ok := config.Deployments[Name]; !ok {
			log.Fatal(fmt.Sprintf("Deployment '%s' not found", Name))
		}
	} else {
		// No deployment name was provided. Find default from config
		if len(config.Deployments) > 1 {
			log.Fatal("Multiple deployments found. Use the --name flag to specify the deployment to use with the command.")
		} else {
			// Use the first and only deployment
			for name, _ := range config.Deployments {
				Name = name
				break
			}
		}
	}

	Project = config.Deployments[Name].Project
}
