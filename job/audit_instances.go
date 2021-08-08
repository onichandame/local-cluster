package job

import (
	"time"

	"github.com/onichandame/local-cluster/instance"
)

var auditInstances = job{
	name:      "AuditInstances",
	immediate: true,
	fatal:     true,
	interval:  time.Minute,
	run: func() error {
		return instance.Audit()
	},
}
