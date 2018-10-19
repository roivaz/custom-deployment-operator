package controller

import (
	"github.com/lominorama/custom-deployment-operator/pkg/controller/customdeployment"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, customdeployment.Add)
}
