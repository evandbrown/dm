package commands

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/util"
	"github.com/nu7hatch/gouuid"
	"github.com/spf13/cobra"
)

var uid bool
var config string
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a configuration to Deployment Manager.",
}

func init() {
	deployCmd.Flags().StringVarP(&config, "config", "c", "config.yaml", "The name of the config to deploy.")
	deployCmd.Flags().BoolVarP(&uid, "uid", "u", true, "Should a 7 char UID be appended to deployment name. Defaults is yes")
	deployCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(deploy(cmd, args))
	}
}

func deploy(cmd *cobra.Command, args []string) error {
	if Project == "" {
		log.WithFields(log.Fields{
			"missingParam": "--project",
		}).Fatal("--project parameter is required to create a new deployment")
	}
	log.Debug("Creating deployment manager service")
	service, err := googlecloud.GetService()
	util.Check(err)

	if len(Name) == 0 {
		Name, err = os.Getwd()
		util.Check(err)
		dirs := strings.Split(Name, "/")
		Name = dirs[len(dirs)-1]
	}

	if uid {
		u, err := uuid.NewV4()
		util.Check(err)
		Name += "-" + u.String()[:7]
	}
	log.Printf("Creating new deployment %s", Name)

	d := googlecloud.NewDeployment(Name, "", config)
	d.Intent = "UPDATE"
	call := service.Deployments.Insert(Project, d)
	_, error := call.Do()
	util.Check(error)
	dConfig := conf.Deployment{
		Id:      Name,
		Project: Project,
	}

	_, err = conf.AppendDeployment(dConfig, true)
	if err != nil {
		log.Fatal(fmt.Sprintf("Config was deployed but there was an error writing the config file. You will not be able to use other `dm` commands, but the deployment will exist. Error was %s", err))
	}

	log.Printf("Deployment created. Run `dm status` for information on its progress")
	return nil
}
