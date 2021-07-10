package job

type job struct {
	name           string
	run            func() error
	fatal          bool
	immediate      bool
	blocking       bool
	interval       string
	successfulRuns uint
	totalRuns      uint
	dependsOn      []*job
}
