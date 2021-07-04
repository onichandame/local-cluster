package job

type job struct {
	name      string
	run       func() error
	fatal     bool
	immediate bool
	interval  string
}
