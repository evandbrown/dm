package commands

import (
	"strings"

	//log "github.com/Sirupsen/logrus"
	"github.com/google/google-api-go-client/deploymentmanager/v2"
)

func getDeploymentStatus(d *deploymentmanager.Deployment) (string, error) {
	var status string
	if d.Operation == nil {
		return status, nil
	}

	status = strings.ToUpper(d.Operation.OperationType) + ": "

	if d.Operation.Error != nil {
		status += "ERROR ("
		for i, e := range d.Operation.Error.Errors {
			status += e.Code
			if i != len(d.Operation.Error.Errors)-1 {
				status += ", "
			}
		}
		status += ")"
	} else {
		status += d.Operation.Status
	}

	return status, nil
}
