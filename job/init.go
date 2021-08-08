package job

import (
	"time"

	"github.com/sirupsen/logrus"
)

func InitJobs() {
	allJobs := []*job{&createAdmin, &runDashboard, &auditInstances, &initInterfaces, &initProxyManager}

	initJob := func(job *job) {
		logrus.Infof("initializing job %s", job.name)
		if job.interval != 0 {
			var runForInterval func()
			runForInterval = func() {
				time.Sleep(job.interval)
				runJob(job)
				runForInterval()
			}
			go runForInterval()
		}
		if job.immediate {
			go runJob(job)
		}
		logrus.Infof("initialized job %s", job.name)
	}

	for _, job := range allJobs {
		initJob(job)
	}
}
