package app

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/onichandame/local-cluster/db/model"
)

func getSpec(appDef *model.Application) (*model.ApplicationSpec, error) {
	var spec *model.ApplicationSpec
	for _, s := range appDef.Specs {
		if s.Platform == runtime.GOOS && s.Arch == runtime.GOARCH {
			spec = &s
		}
	}
	if spec == nil {
		return nil, errors.New(fmt.Sprintf("failed to find the spec for the runtime! Platform = %s Arch = %s", runtime.GOOS, runtime.GOARCH))
	}
	return spec, nil

}
