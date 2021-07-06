package job

import (
	"errors"

	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var createAdmin = job{
	interval:       "1m",
	name:           "CreateAdmin",
	fatal:          true,
	immediate:      true,
	successfulRuns: 1,
	run: func() error {
		var err error
		admin := model.User{}
		if err = db.Db.Where("role = ?", model.ADMIN).First(&admin).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = db.Db.Transaction(func(db *gorm.DB) error {
					const keyAdminName = "admin_name"
					viper.SetDefault(keyAdminName, "admin")
					admin.Name = viper.GetString(keyAdminName)
					admin.Role = model.ADMIN
					err := db.Create(&admin).Error
					if err != nil {
						return err
					}
					const keyAdminPass = "admin_password"
					viper.SetDefault(keyAdminPass, "admin")
					cred := model.Credential{UserID: admin.ID, Password: viper.GetString(keyAdminPass)}
					err = db.Create(&cred).Error
					if err != nil {
						return err
					}
					return nil
				})
			}
		}
		return err
	},
}
