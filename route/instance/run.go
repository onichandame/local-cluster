package instance

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/constants"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance"
	"github.com/onichandame/local-cluster/pkg/route"
)

var InstanceRun = route.Route{Endpoint: "/run", Method: route.POST, Handler: func(c *gin.Context) (interface{}, error) {
	type RequestBody struct {
		ApplicationID uint                            `json:"application_id" binding:"required"`
		RestartPolicy constants.InstanceRestartPolicy `json:"restart_policy" binding:"required"`
	}
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		return nil, err
	}
	ins := model.Instance{}
	ins.ApplicationID = requestBody.ApplicationID
	ins.RestartPolicy = requestBody.RestartPolicy
	if !constants.IsValidInstanceRestartPolicy(ins.RestartPolicy) {
		return nil, errors.New(fmt.Sprintf("restart policy %s not recognized!", ins.RestartPolicy))
	}
	if err := instance.RunInstance(&ins); err != nil {
		return nil, err
	}
	return &ins, nil
}}
