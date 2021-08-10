package job

import (
	"time"

	"github.com/onichandame/local-cluster/template"
)

var auditTemplates = job{
	name: "AuditTemplates",
	run: func() (err error) {
		err = template.Audit()
		return err
	},
	fatal:     true,
	immediate: true,
	interval:  time.Minute * 5,
}
