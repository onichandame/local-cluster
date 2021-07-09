package app

import "github.com/onichandame/local-cluster/pkg/route"

var AppRoot = route.Route{
	Subroutes: []*route.Route{&AppAdd, &AppList},
	Endpoint:  "/apps",
}
