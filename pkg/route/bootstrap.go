package route

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Bootstrap(r *Route, g *gin.RouterGroup) {
	group := g.Group(r.Endpoint)
	handleRequest := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			res, err := r.Handler(c)
			if err != nil {
				code := 400
				c.JSON(code, gin.H{"message": err.Error()})
				return
			}
			var body []byte
			var contentType string
			if r, ok := res.([]byte); ok {
				body = r
				contentType = "text/plain"
			} else if r, ok := res.(string); ok {
				body = []byte(r)
				contentType = "text/plain"
			} else {
				logrus.Info(res)
				if body, err = json.Marshal(res); err != nil {
					c.JSON(500, gin.H{"message": err.Error()})
					return
				}
				logrus.Info(body)
				contentType = "application/json"
			}
			c.Data(200, contentType, body)
		}
	}
	if r.Handler != nil {
		switch r.Method {
		case GET:
			group.GET("", handleRequest())
		case POST:
			group.POST("", handleRequest())
		case PUT:
			group.PUT("", handleRequest())
		case DELETE:
			group.DELETE("", handleRequest())
		case ANY:
			fallthrough
		default:
			group.Any("", handleRequest())
		}
	}
	for _, sub := range r.Subroutes {
		Bootstrap(sub, group)
	}
}
