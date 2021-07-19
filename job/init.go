package job

import (
	"time"

	"github.com/chebyrash/promise"
	"github.com/sirupsen/logrus"
)

var jobInitMap = make(map[*job]*promise.Promise)

func JobsInit() {
	allJobs := []*job{&createAdmin, &runDashboard, &auditInstances, &initInterfaces, &initProxyManager, &initCacheManager}
	for _, j := range allJobs {
		if _, ok := jobInitMap[j]; !ok {
			initAJob(j)
		}
	}
	allPs := []*promise.Promise{}
	for _, p := range jobInitMap {
		allPs = append(allPs, p)
	}
	if _, err := promise.All(allPs...).Await(); err != nil {
		logrus.Error(err)
		panic("failed to init all jobs")
	}
}

func initAJob(j *job) {
	// skip if already inited
	if _, ok := jobInitMap[j]; ok {
		return
	}
	logrus.Infof("initializing job %s", j.name)
	initInterval := func() {
		duration, err := time.ParseDuration(j.interval)
		if err != nil {
			logrus.Fatalf("job %s failed initialization!", j.name)
		}
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				go runJob(j)
			}
		}
	}
	for _, dep := range j.dependsOn {
		initAJob(dep)
	}
	jobInitMap[j] = promise.New(func(resolve func(promise.Any), reject func(error)) {
		depPs := []*promise.Promise{}
		for _, dep := range j.dependsOn {
			depPs = append(depPs, jobInitMap[dep])
		}
		if _, err := promise.All(depPs...).Await(); err != nil {
			reject(err)
			return
		}
		if j.interval != "" {
			go initInterval()
		}
		if j.immediate {
			if j.blocking {
				if err := runJob(j); err != nil {
					panic(err)
				}
			} else {
				go runJob(j)
			}
		}
		logrus.Infof("initialized job %s", j.name)
		resolve(nil)
	})
}
