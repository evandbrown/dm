package main

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

var path = ".dm.toml"

type DeploymentConfig struct {
	Deployments map[string]Deployment
}

type Deployment struct {
	Id          string
	Project     string
	Config      string
	Vars        string
	Fingerprint string
	Commit      string
	Deployed    string
	Updated     string
	Branch      string
}

func ReadDeploymentConfig() (DeploymentConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return DeploymentConfig{}, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return DeploymentConfig{}, err
	}
	var config DeploymentConfig
	if _, err := toml.Decode(string(buf), &config); err != nil {
		return DeploymentConfig{}, err
	}
	return config, nil
}

func AppendDeployment(d Deployment, initIfMissing bool) (DeploymentConfig, error) {
	f, err := os.Open(path)
	if err != nil && initIfMissing {
		f, err = os.Create(path)
	}
	f.Close()

	config, err := ReadDeploymentConfig()
	if err != nil {
		return DeploymentConfig{}, err
	}
	if config.Deployments == nil {
		config.Deployments = make(map[string]Deployment)
	}

	config.Deployments[d.Id] = d

	f, _ = os.Create(path)
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return DeploymentConfig{}, err
	}
	return config, nil
}

func RemoveDeployment(id string) error {
	config, err := ReadDeploymentConfig()
	if err != nil {
		return err
	}

	if _, ok := config.Deployments[id]; ok == false {
		return errors.New("Deployment not found")
	}

	delete(config.Deployments, id)
	if len(config.Deployments) == 0 {
		return os.Remove(path)
	}

	f, _ := os.Create(path)
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}
