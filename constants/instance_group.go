package constants

type InstanceGroupStatus string

const (
	INITIALIZING InstanceGroupStatus = "INITIALIZING"
	READY        InstanceGroupStatus = "READY"
	NOTREADY     InstanceGroupStatus = "NOTREADY"
	SHUTDOWN     InstanceGroupStatus = "SHUTDOWN"
)
