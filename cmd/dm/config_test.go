package main

import (
	"testing"
)

func TestReadDeploymentConfig(t *testing.T) {
	_, err := ReadDeploymentConfig()
	if err == nil {
		t.Error("Expected error when reading deployment configuration")
	}
}

func TestAppendDeployment(t *testing.T) {
	d1 := Deployment{Id: "test1", Project: "testproject1"}
	dc, err := AppendDeployment(d1, true)
	if err != nil {
		t.Error("Error creating a new deployment")
	}
	if len(dc.Deployments) != 1 {
		t.Error("Expected deployment config length to be 1")
	}

	d2 := Deployment{Id: "test2", Project: "testproject2"}
	dc, err = AppendDeployment(d2, true)
	if err != nil {
		t.Error("Error creating a new deployment")
	}
	if len(dc.Deployments) != 2 {
		t.Error("Expected deployment config length to be 2")
	}

	RemoveDeployment(d1.Id)
	dc, err = ReadDeploymentConfig()
	if err != nil {
		t.Error("Error reading deployment configuration")
	}
	if len(dc.Deployments) != 1 {
		t.Error("Expected deployment config length to be 1")
	}

	RemoveDeployment(d2.Id)
	dc, err = ReadDeploymentConfig()
	if err == nil {
		t.Error("Deployment config file should not exist")
	}
}
