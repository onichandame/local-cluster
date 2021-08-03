package gateway

type GatewayStatus string

// (create) → CREATING → READY →
// 											   								↑			↓   TERMINATING → TERMINATED → (delete)
//                    NOTREADY →

const (
	CREATING    GatewayStatus = "CREATING"
	READY       GatewayStatus = "CREATING"
	NOTREADY    GatewayStatus = "NOTREADY"
	TERMINATING GatewayStatus = "TERMINATING"
	TERMINATED  GatewayStatus = "TERMINATED"
)
