package interfaces

import "github.com/onichandame/local-cluster/db/model"

func deleteIF(svcIf *model.ServiceInterface) error {
	delete(LockedPortsMap, svcIf.Port)
	return nil
}
