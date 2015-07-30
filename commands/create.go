package commands

import (
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/util"
	"github.com/nu7hatch/gouuid"
	"github.com/spf13/cobra"
)

var uid bool
var config, name string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new deployment",
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
	//TODO validate required params (i.e., project)
	log.Debug("Creating deployment manager service")
	service, err := util.GetService()
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

	d := util.NewDeployment(name, "", config)
	d.Intent = "UPDATE"
	call := service.Deployments.Insert(cmd.Flags().Lookup("project").Value.String(), d)
	_, error := call.Do()
	util.Check(error)
	log.Printf("Deployment created. Run `dm status` for information on its progress")
	return nil
}
