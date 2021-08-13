package job

import (
	"time"

	"github.com/onichandame/local-cluster/instance_group"
)

var auditInstanceGroups = job{
	name:      "AuditInstanceGroups",
	immediate: true,
	fatal:     true,
	interval:  time.Minute * 5,
	run: func() error {
		return instancegroup.Audit()
	},
}
