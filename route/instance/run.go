package instance

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/instance"
	"github.com/onichandame/local-cluster/pkg/route"
)

var InstanceRun = route.Route{Endpoint: "/run", Method: route.POST, Handler: func(c *gin.Context) (interface{}, error) {
	type RequestBody struct {
		ApplicationID uint            `json:"application_id" binding:"required"`
		RestartPolicy model.EnumValue `json:"restart_policy" binding:"required"`
	}
	var requestBody RequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		return nil, err
	}
	ins := model.Instance{}
	ins.ApplicationID = requestBody.ApplicationID
	policies := model.GetRestartPolicies(db.Db)
	policy, ok := policies[requestBody.RestartPolicy]
	if !ok {
		return nil, errors.New(fmt.Sprintf("restart policy %s not recognized!", requestBody.RestartPolicy))
	}
	ins.RestartPolicyID = policy.ID
	if err := instance.RunInstance(&ins); err != nil {
		return nil, err
	}
	return &ins, nil
}}
