package job

import (
	"time"

	"github.com/onichandame/local-cluster/instance"
)

var auditInstances = job{
	name:      "AuditInstances",
	immediate: true,
	fatal:     true,
	interval:  time.Minute * 5,
	run: func() error {
		return instance.Audit()
	},
}
