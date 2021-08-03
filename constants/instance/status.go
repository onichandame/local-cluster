package instance

type InstanceStatus string

// (create) → CREATING → RUNNING → TERMINATING → TERMINATED → (delete)
// 											   								↑					↓
//               RESTARTING	←	CRASHED → (delete)

const (
	CREATING    InstanceStatus = "CREATING"
	RUNNING     InstanceStatus = "RUNNING"
	CRASHED     InstanceStatus = "CRASHED"
	RESTARTING  InstanceStatus = "RESTARTING"
	TERMINATING InstanceStatus = "TERMINATING"
	TERMINATED  InstanceStatus = "TERMINATED"
)
