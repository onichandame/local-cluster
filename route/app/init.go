package app

import "github.com/onichandame/local-cluster/pkg/route"

var AppInit = route.Route{Subroutes: []*route.Route{&AppRun}, Endpoint: "/app"}
