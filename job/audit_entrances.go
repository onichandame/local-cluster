package job

import "github.com/onichandame/local-cluster/entrance"

var auditEntrances = job{
	name:      "AuditEntrances",
	immediate: true,
	fatal:     true,
	blocking:  true,
	dependsOn: []*job{&auditInstances},
	run: func() error {
		return entrance.Audit()
	},
}
