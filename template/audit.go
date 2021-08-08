package template

import (
	"errors"

	"github.com/chebyrash/promise"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"gorm.io/gorm"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.Template{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		proms := []*promise.Promise{}
		for rows.Next() {
			var template model.Template
			if err = db.Db.ScanRows(rows, &template); err != nil {
				panic(err)
			}
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
				var err error
				var ins *model.Instance
				if err = db.Db.First(ins, "template_name = ?", template.Name).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						panic(err)
					}
				}
				var ig *model.InstanceGroup
				if err = db.Db.First(ig, "template_name = ?", template.Name).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						panic(err)
					}
				}
				if ins == nil && ig == nil {
					if err = db.Db.Delete(&template).Error; err != nil {
						panic(err)
					}
				}
			}))
		}
		if _, err = promise.AllSettled(proms...).Await(); err != nil {
			panic(err)
		}
	}
	if _, err = promise.AllSettled([]*promise.Promise{
		promise.New(func(resolve func(promise.Any), reject func(error)) {
			defer utils.SettlePromise(resolve, reject)
			if err := auditStorageBindings(); err != nil {
				panic(err)
			}
		}), promise.New(func(resolve func(promise.Any), reject func(error)) {
			defer utils.SettlePromise(resolve, reject)
			if err := auditProbes(); err != nil {
				panic(err)
			}
		}),
	}...).Await(); err != nil {
		panic(err)
	}
	return err
}

func auditStorageBindings() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.StorageBinding{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		proms := []*promise.Promise{}
		for rows.Next() {
			var sb model.StorageBinding
			if err = db.Db.ScanRows(rows, &sb); err != nil {
				panic(err)
			}
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
				var err error
				if err = db.Db.First(&model.Template{}, sb.TemplateID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = db.Db.Delete(&sb).Error; err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}
			}))
		}
		if _, err = promise.AllSettled(proms...).Await(); err != nil {
			panic(err)
		}
	}
	return err
}

func auditProbes() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.Probe{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		proms := []*promise.Promise{}
		for rows.Next() {
			var probe model.Probe
			if err = db.Db.ScanRows(rows, &probe); err != nil {
				panic(err)
			}
			proms = append(proms, promise.New(func(resolve func(promise.Any), reject func(error)) {
				defer utils.SettlePromise(resolve, reject)
				var err error
				if err = db.Db.First(&model.Template{}, probe.TemplateID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = db.Db.Delete(&probe).Error; err != nil {
							panic(err)
						}
						db.Db.Where("probe_id = ?", probe.ID).Delete(&model.TCPProbe{})
						db.Db.Where("probe_id = ?", probe.ID).Delete(&model.HTTPProbe{})
					}
				}
			}))
		}
		if promise.All(proms...).Await(); err != nil {
			panic(err)
		}
	}
	return err
}
