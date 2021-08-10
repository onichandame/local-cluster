package template

import (
	"errors"
	"sync"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Audit() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.Template{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		wg := sync.WaitGroup{}
		for rows.Next() {
			wg.Add(1)
			var template model.Template
			if err = db.Db.Preload("Instances").Preload("InstanceGroups").ScanRows(rows, &template); err != nil {
				panic(err)
			}
			go func() {
				defer wg.Done()
				var err error
				defer func() {
					if err != nil {
						logrus.Error(err)
					}
				}()
				defer utils.RecoverFromError(&err)
				if len(template.Instances) == 0 && len(template.InstanceGroups) == 0 {
					if err = db.Db.Delete(&template).Error; err != nil {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := auditStorageBindings(); err != nil {
			panic(err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := auditProbes(); err != nil {
			panic(err)
		}
	}()
	wg.Wait()
	return err
}

func auditStorageBindings() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.StorageBinding{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		wg := sync.WaitGroup{}
		for rows.Next() {
			var sb model.StorageBinding
			if err = db.Db.ScanRows(rows, &sb); err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				var err error
				defer func() {
					if err != nil {
						logrus.Error(err)
					}
				}()
				defer utils.RecoverFromError(&err)
				if err = db.Db.First(&model.Template{}, sb.TemplateID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = db.Db.Delete(&sb).Error; err != nil {
							panic(err)
						}
					} else {
						panic(err)
					}
				}
			}()
		}
		wg.Wait()
	}
	return err
}

func auditProbes() (err error) {
	defer utils.RecoverFromError(&err)
	if rows, err := db.Db.Model(&model.Probe{}).Rows(); err != nil {
		panic(err)
	} else {
		defer rows.Close()
		wg := sync.WaitGroup{}
		for rows.Next() {
			var probe model.Probe
			if err = db.Db.ScanRows(rows, &probe); err != nil {
				panic(err)
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				var err error
				defer func() {
					if err != nil {
						logrus.Error(err)
					}
				}()
				defer utils.RecoverFromError(&err)
				if err = db.Db.First(&model.Template{}, probe.TemplateID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						if err = db.Db.Delete(&probe).Error; err != nil {
							panic(err)
						}
						db.Db.Where("probe_id = ?", probe.ID).Delete(&model.TCPProbe{})
						db.Db.Where("probe_id = ?", probe.ID).Delete(&model.HTTPProbe{})
					}
				}
			}()
		}
		wg.Wait()
	}
	return err
}
