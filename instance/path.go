package instance

import (
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
)

func getRuntimeDir(runtime string) string {
	return filepath.Join(config.Config.Path.Instances, runtime)
}

func getInsRuntimeDir(id uint) string {
	return getRuntimeDir(strconv.Itoa(int(id)))
}
