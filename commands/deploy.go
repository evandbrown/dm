package commands

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/conf"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/evandbrown/dm/util"
	"github.com/nu7hatch/gouuid"
	"github.com/spf13/cobra"
)

const (
	uidlen = 5
	maxlen = 63
	namere = "[a-z]([-a-z0-9]*[a-z0-9])?"
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

func getName(setUid bool) (string, error) {
	var name string
	var err error
	if len(Name) == 0 {
		name, err = os.Getwd()
		if err != nil {
			return "", err
		}
		dirs := strings.Split(name, "/")
		Name = dirs[len(dirs)-1]
	} else {
		name = Name
	}

	// Replace underscores
	name = strings.Replace(name, "_", "-", -1)
	name = strings.ToLower(name)

	// Reduce name prefix to keep total to < 63 chars
	if setUid && len(name)+uidlen > maxlen {
		name = name[:maxlen-uidlen]
	}

	// Append a uid
	if setUid {
		u, err := uuid.NewV4()
		util.Check(err)
		name += "-" + u.String()[:uidlen-1]
	}

	// Validate name
	if match, err := regexp.MatchString(namere, name); match == false || err != nil {
		return "", errors.New(fmt.Sprintf("The provided or derived name for the deployment is invalid: %s. Must match regex %s", name, namere))
	}
	return name, nil
}

func deploy(cmd *cobra.Command, args []string) error {
	if Project == "" {
		log.Fatal("--project parameter is required to create a new deployment")
	}
	log.Debug("Creating deployment manager service")
	service, err := googlecloud.GetService()
	util.Check(err)

	Name, err = getName(uid)
	if err != nil {
		log.Warning(err)
		return err
	}

	log.Infof("Creating new deployment %s", Name)

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

	fmt.Printf("Created deployment %s.\n", Name)
	return nil
}
