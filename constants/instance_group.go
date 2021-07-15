package constants

type InstanceGroupStatus string
type UpdateStrategy string

const (
	INITIALIZING InstanceGroupStatus = "INITIALIZING"
	READY        InstanceGroupStatus = "READY"
	NOTREADY     InstanceGroupStatus = "NOTREADY"
	SHUTDOWN     InstanceGroupStatus = "SHUTDOWN"
)
