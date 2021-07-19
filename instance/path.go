package instance

import (
	"path/filepath"
	"strconv"

	"github.com/onichandame/local-cluster/config"
	"github.com/onichandame/local-cluster/db/model"
)

func getInsDir(insDef *model.Instance) string {
	return filepath.Join(config.Config.Path.Instances, strconv.Itoa(int(insDef.ID)))
}
