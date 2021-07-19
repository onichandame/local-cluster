package job

import "github.com/onichandame/local-cluster/instance"

var auditInstances = job{
	name:      "AuditInstances",
	immediate: true,
	blocking:  true,
	fatal:     true,
	run: func() error {
		return instance.Audit()
	},
}
