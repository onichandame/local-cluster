package job

import "time"

type job struct {
	name           string
	run            func() (err error)
	fatal          bool // whether panic if job runs and fails
	immediate      bool // whether run once right after main application startup
	interval       time.Duration
	successfulRuns uint
	totalRuns      uint
}
