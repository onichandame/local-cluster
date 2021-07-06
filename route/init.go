package route

import (
	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/pkg/route"
)

func RoutesInit(e *gin.Engine) {
	root := e.Group("/")

	rootRoute := route.Route{
		Endpoint:  "",
		Subroutes: []*route.Route{&RouteHealthCheck},
	}

	route.Bootstrap(&rootRoute, root)
}
