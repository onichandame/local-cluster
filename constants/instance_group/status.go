package instancegroup

type InstanceGroupStatus string

// (create) → CREATING → NOTREADY →
// 																						   ↑   ↓   TERMINATING → TERMINATED → (delete)
//                          READY →
const (
	CREATING    InstanceGroupStatus = "CREATING"
	READY       InstanceGroupStatus = "READY"
	NOTREADY    InstanceGroupStatus = "NOTREADY"
	TERMINATING InstanceGroupStatus = "TERMINATING"
	TERMINATED  InstanceGroupStatus = "TERMINATED"
)
