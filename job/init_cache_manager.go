package job

import "github.com/onichandame/local-cluster/application"

var initCacheManager = job{
	name:      "InitCacheManager",
	immediate: true,
	blocking:  true,
	dependsOn: []*job{&initConfig},
	run: func() error {
		return application.InitCache()
	},
}
