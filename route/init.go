package route

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/pkg/route"
	"github.com/onichandame/local-cluster/route/app"
	"github.com/onichandame/local-cluster/route/instance"
)

func RoutesInit(e *gin.Engine) {
	root := e.Group("/")

	rootRoute := route.Route{
		Endpoint:  "",
		Subroutes: []*route.Route{&RouteHealthCheck, &app.AppRoot, &instance.InstanceRoot},
	}

	route.Bootstrap(&rootRoute, root)
}
