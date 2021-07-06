package app

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/route"
)

var AppAdd = route.Route{Endpoint: "/add", Method: route.POST, Handler: func(c *gin.Context) (interface{}, error) {
	type RequestBody struct {
		Name  string `json:"name" binding:"required"`
		Specs []struct {
			Platform    string `json:"platform" binding:"requred"`
			Arch        string `json:"arch" binding:"requred"`
			Target      string `json:"target" binding:"required"`
			DownloadUrl string `json:"download_url" binding:"required"`
			Hash        string `json:"hash"`
		} `json:"specs" binding:"required"`
	}
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		return nil, err
	}
	app := model.Application{}
	app.Name = requestBody.Name
	app.Specs = requestBody.Specs
	db.Db.Create
	return nil, nil
}}
