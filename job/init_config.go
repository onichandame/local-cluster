package job

import "github.com/onichandame/local-cluster/config"

var initConfig = job{
	name:      "InitConfig",
	immediate: true,
	fatal:     true,
	blocking:  true,
	run: func() error {
		return config.Init()
	},
}
