package constants

type InstanceStatus string
type InstanceRestartPolicy string

const (
	CREATING    InstanceStatus = "CREATING"
	RUNNING     InstanceStatus = "RUNNING"
	TERMINATING InstanceStatus = "TERMINATING"
	CRASHED     InstanceStatus = "CRASHED"
	TERMINATED  InstanceStatus = "TERMINATED"

	ALWAYS    InstanceRestartPolicy = "ALWAYS"
	NEVER     InstanceRestartPolicy = "NEVER"
	ONFAILURE InstanceRestartPolicy = "ONFAILURE"
)

func IsValidInstanceRestartPolicy(raw InstanceRestartPolicy) bool {
	mapper := map[InstanceRestartPolicy]interface{}{ALWAYS: nil, NEVER: nil, ONFAILURE: nil}
	_, ok := mapper[raw]
	return ok
}
