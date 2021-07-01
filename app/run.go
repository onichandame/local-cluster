package app

import (
	"os/exec"

	"github.com/onichandame/local-cluster/pkg/types"
)

func AppRun(appDef *types.AppDefinition) {
	command, err := exec.LookPath(appDef.Command)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(command, appDef.Entry)
}
