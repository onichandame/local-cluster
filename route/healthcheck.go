package route

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/pkg/route"
)

var RouteHealthCheck = route.Route{
	Endpoint: "/healthcheck",
	Method:   route.GET,
	Handler: func(c *gin.Context) (interface{}, error) {
		return map[string]interface{}{"timestamp": time.Now().UnixNano() / int64(time.Millisecond),
			"ok": true,
		}, nil
	},
}
