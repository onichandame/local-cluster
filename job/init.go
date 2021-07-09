package job

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func JobInit() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go InitAJob(&createAdmin, &wg)
	wg.Add(1)
	go InitAJob(&runDashboard, &wg)
	logrus.Info("hi")

	wg.Wait()
}

func InitAJob(job *job, wg *sync.WaitGroup) {
	logrus.Infof("initializing job %s", job.name)
	defer wg.Done()
	initInterval := func() {
		duration, err := time.ParseDuration(job.interval)
		if err != nil {
			logrus.Fatalf("job %s failed initialization!", job.name)
		}
		ticker := time.NewTicker(duration)
		for {
			select {
			case <-ticker.C:
				go runJob(job)
			}
		}
	}
	if job.interval != "" {
		go initInterval()
	}
	if job.immediate {
		runJob(job)
	}
	logrus.Infof("initialized job %s", job.name)
}
