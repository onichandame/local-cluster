package instance

type InstanceStatus string

// (create) → CREATING → WAITING → RUNNING → TERMINATING → TERMINATED → (delete)
// 											             				↓			↓    ↑
//                            CRASHED → RESTARTING
// ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓	↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓	↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓
// 																														FAILED

const (
	CREATING    InstanceStatus = "CREATING"
	WAITING     InstanceStatus = "WAITING"
	RUNNING     InstanceStatus = "RUNNING"
	CRASHED     InstanceStatus = "CRASHED"
	RESTARTING  InstanceStatus = "RESTARTING"
	TERMINATING InstanceStatus = "TERMINATING"
	TERMINATED  InstanceStatus = "TERMINATED"
	FAILED      InstanceStatus = "FAILED"
)
