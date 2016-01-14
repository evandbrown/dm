package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/evandbrown/dm/googlecloud"
	"github.com/nu7hatch/gouuid"
	"github.com/spf13/cobra"
)

const (
	uidlen            = 5
	maxlen            = 63
	namere            = "[a-z]([-a-z0-9]*[a-z0-9])?"
	configpathDefault = "config.yaml"
	varspathDefault   = "vars.yaml"
)

var uid bool
var configpath, varspath string
var vars configVar
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a configuration to Deployment Manager.",
}

func init() {
	deployCmd.Flags().StringVarP(&configpath, "config-file", "c", configpathDefault, "The name of the config to deploy.")
	deployCmd.Flags().VarP(&vars, "var", "v", "A variable value to provide to the vars.yaml file for use in a deployment. Define multiple with -v var1=foo -v var2=2")
	deployCmd.Flags().StringVarP(&varspath, "vars-file", "x", varspathDefault, "The name of the config to deploy.")
	deployCmd.Flags().BoolVarP(&uid, "uid", "u", true, "Should a 7 char UID be appended to deployment name. Defaults is yes")
	deployCmd.Run = func(cmd *cobra.Command, args []string) {
		Check(deploy(cmd, args))
	}
}

// Special flag type to accumulate vars
type configVar struct {
	vars map[string]string
}

// Implemenet the flag interface
func (v *configVar) String() string {
	return fmt.Sprint(*v)
}

// Implemenet the flag interface
func (v *configVar) Set(value string) error {
	if v.vars == nil {
		v.vars = make(map[string]string)
	}
	s := strings.Split(value, "=")
	if len(s) != 2 {
		return errors.New("value must be formatted as k=v")
	}
	v.vars[strings.TrimSpace(s[0])] = strings.TrimSpace(s[1])
	return nil
}

func (v *configVar) Type() string {
	return "Convig var"
}

func deploy(cmd *cobra.Command, args []string) error {
	//TODO this should be a validation method
	if Project == "" {
		log.Fatal("--project parameter is required to create a new deployment")
	}

	service, err := googlecloud.GetService()
	Check(err)

	Name, err = getName(uid)
	if err != nil {
		log.Warning(err)
		return err
	}

	log.Infof("Creating new deployment %s", Name)
	depBuilder := &DeploymentBuilder{
		DeploymentName:  Name,
		DeploymentDesc:  "",
		ConfigFilePath:  configpath,
		VarsDotYamlPath: varspath,
		CLIVars:         vars.vars,
	}

	d, err := depBuilder.GetDeployment()
	if err != nil {
		log.Warning(err)
		return err
	}

	//	d.Intent = "UPDATE"
	call := service.Deployments.Insert(Project, d)
	_, error := call.Do()
	Check(error)

	//TODO only set Vars if the varspath file actually exists
	dConfig := Deployment{
		Id:      Name,
		Project: Project,
		Config:  configpath,
		Vars:    varspath,
	}

	_, err = AppendDeployment(dConfig, true)
	if err != nil {
		log.Fatal(fmt.Sprintf("Config was deployed but there was an error writing the config file. You will not be able to use other `dm` commands, but the deployment will exist. Error was %s", err))
	}

	fmt.Printf("Created deployment %s.\n", Name)
	return nil
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
		name = dirs[len(dirs)-1]
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
		Check(err)
		name += "-" + u.String()[:uidlen-1]
	}

	// Validate name
	if match, err := regexp.MatchString(namere, name); match == false || err != nil {
		return "", errors.New(fmt.Sprintf("The provided or derived name for the deployment is invalid: %s. Must match regex %s", name, namere))
	}
	return name, nil
}
