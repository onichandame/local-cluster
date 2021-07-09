package app

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/route"
)

var AppList = route.Route{
	Endpoint: "/list",
	Method:   route.POST,
	Handler: func(c *gin.Context) (interface{}, error) {
		res := []*model.Application{}
		if err := db.Db.Find(&res).Error; err != nil {
			return nil, err
		}
		return res, nil
	},
}
