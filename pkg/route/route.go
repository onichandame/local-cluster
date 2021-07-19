package route

import "github.com/gin-gonic/gin"

type method uint

const (
	GET    = iota
	POST   = iota
	PUT    = iota
	DELETE = iota
	ANY    = iota
)

type Route struct {
	Endpoint  string
	Method    method
	Handler   func(c *gin.Context) (interface{}, error)
	Subroutes []*Route
}
