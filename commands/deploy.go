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
var config, name string
var createCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a configuration to Deployment Manager.",
}

func init() {
	createCmd.Flags().StringVarP(&config, "config", "c", "config.yaml", "The name of the config to deploy.")
	createCmd.Flags().StringVarP(&name, "name", "n", "", "The name of the deployment. Defaults to the name of the currect directory")
	createCmd.Flags().BoolVarP(&uid, "uid", "u", true, "Should a 7 char UID be appended to deployment name. Defaults is yes")
	createCmd.Run = func(cmd *cobra.Command, args []string) {
		util.Check(create(cmd, args))
	}
}

func create(cmd *cobra.Command, args []string) error {
	if Project == "" {
		log.WithFields(log.Fields{
			"missingParam": "--project",
		}).Fatal("--project parameter is required to create a new deployment")
	}
	log.Debug("Creating deployment manager service")
	service, err := googlecloud.GetService()
	if err != nil {
		return err
	}

	if len(name) == 0 {
		name, err = os.Getwd()
		util.Check(err)
		dirs := strings.Split(name, "/")
		name = dirs[len(dirs)-1]
	}

	if uid {
		u, err := uuid.NewV4()
		util.Check(err)
		name += "-" + u.String()[:7]
	}
	log.Printf("Creating new deployment %s", name)

	d := googlecloud.NewDeployment(name, "", config)
	d.Intent = "UPDATE"
	call := service.Deployments.Insert(Project, d)
	_, error := call.Do()
	util.Check(error)
	dConfig := conf.Deployment{
		Id:      name,
		Project: Project,
	}

	_, err = conf.AppendOrUpdateDeployment(dConfig, true)
	if err != nil {
		log.Fatal(fmt.Sprintf("Config was deployed but there was an error writing the config file. You will not be able to use other `dm` commands, but the deployment will exist. Error was %s", err))
	}

	log.Printf("Deployment created. Run `dm status` for information on its progress")
	return nil
}
