package app

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/route"
)

var AppRun = route.Route{
	Endpoint: "/run",
	Method:   route.POST,
	Handler: func(c *gin.Context) (interface{}, error) {
		type RequestBody struct {
			ApplicationID uint `form:"application_id" json:"application_id" xml:"application_id" binding:"required"`
		}
		var requestBody RequestBody
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			return nil, err
		}
		app := model.Application{}
		if err := db.Db.First(&app, requestBody.ApplicationID).Error; err != nil {
			return nil, err
		}
		instance, err := initInstance(&app)
		if err != nil {
			return nil, err
		}
		preapreInstance(instance)
		return nil, nil
	},
}

func initInstance(app *model.Application) (*model.Instance, error) {
	statuses := model.GetInstanceStatuses(db.Db)
	instance := model.Instance{ApplicationID: app.ID, StatusID: statuses[model.PENDING].ID}
	return &instance, db.Db.Create(&instance).Error
}

func preapreInstance(instance *model.Instance) {
}
