package job

type job struct {
	name           string
	run            func() error
	fatal          bool // whether panic if job runs and fails
	immediate      bool // whether run once right after main application startup
	blocking       bool // during immediate run, whether blocks the main thread
	interval       string
	successfulRuns uint
	totalRuns      uint
	dependsOn      []*job // only start initializing after the dependencies have been initialized
}
