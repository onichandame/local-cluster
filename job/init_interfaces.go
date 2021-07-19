package job

import (
	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/interfaces"
)

var initInterfaces = job{
	name:      "InitInterfaces",
	immediate: true,
	fatal:     true,
	blocking:  true,
	run: func() error {
		return interfaces.Init(config.Config.PortRange.StartAt, config.Config.PortRange.EndAt)
	},
}
