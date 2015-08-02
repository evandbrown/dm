package conf

import (
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

var path = ".dm.toml"

type DeploymentConfig struct {
	Deployment []Deployment
}

type Deployment struct {
	Id       string
	Project  string
	Commit   string
	Deployed string
	Updated  string
	Branch   string
}

func ReadDeploymentConfig(path string) (DeploymentConfig, error) {
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
	if err := toml.Unmarshal(buf, &config); err != nil {
		return DeploymentConfig{}, err
	}
	return config, nil
}

func AppendOrUpdateDeployment(d Deployment, initIfMissing bool) (DeploymentConfig, error) {
	f, err := os.Open(path)
	if err != nil && initIfMissing {
		f, err = os.Create(path)
	}
	defer f.Close()

	config, err := ReadDeploymentConfig(path)
	if err != nil {
		return DeploymentConfig{}, err
	}
	if config.Deployment == nil {
		config.Deployment = make([]Deployment, 0)
	}

	config.Deployment = append(config.Deployment, d)
	buff, err := toml.Marshal(config)
	if err != nil {
		return DeploymentConfig{}, err
	}
	_, err = f.Write(buff)
	if err != nil {
		return DeploymentConfig{}, err
	}
	return config, nil
}
