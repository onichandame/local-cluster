package job

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var jobInitMap = make(map[*job]*sync.WaitGroup)
var initedWG = sync.WaitGroup{}

func JobInit() {
	initedWG.Add(1)
	allJobs := []*job{&createAdmin, &runDashboard, &auditInstances}
	for _, j := range allJobs {
		jobInitMap[j] = initAJob(j)
	}
	initedWG.Done()
}

func initAJob(j *job) *sync.WaitGroup {
	logrus.Infof("initializing job %s", j.name)
	wg := sync.WaitGroup{}
	wg.Add(1)
	defer func() {
		if !j.blocking {
			wg.Done()
		}
	}()
	initInterval := func() {
		duration, err := time.ParseDuration(j.interval)
		if err != nil {
			logrus.Fatalf("job %s failed initialization!", j.name)
		}
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				go runJob(j, nil)
			}
		}
	}
	if j.interval != "" {
		go initInterval()
	}
	if j.immediate {
		if j.blocking {
			go runJob(j, &wg)
		} else {
			go runJob(j, nil)
		}
	}
	logrus.Infof("initialized job %s", j.name)
	return &wg
}
