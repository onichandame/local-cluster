package config

type config struct {
	Path struct {
		Root      string
		DB        string
		Cache     string
		Instances string
		Storage   string
	}
	PortRange struct {
		StartAt uint
		EndAt   uint
	}
}

func newConfig() *config {
	return new(config)
}

var Config *config
