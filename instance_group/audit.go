package instancegroup

import (
	"sync"

	igConstants "github.com/onichandame/local-cluster/constants/instance_group"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	// audit ig in db
	if rows, err := db.Db.Model(&model.Instance{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		var wg sync.WaitGroup
		for rows.Next() {
			var ig model.InstanceGroup
			if err := db.Db.Preload("Template").ScanRows(rows, &ig); err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
			}()
		}
		wg.Wait()
	}
	return err
}
