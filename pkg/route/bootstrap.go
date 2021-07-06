package route

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
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
			if res, ok := res.([]byte); ok {
				body = res
				contentType = "text/plain"
			}
			if res, ok := res.(string); ok {
				body = []byte(res)
				contentType = "text/plain"
			}
			if res, ok := res.(map[string]interface{}); ok {
				if body, err = json.Marshal(res); err != nil {
					c.JSON(500, gin.H{"message": err.Error()})
					return
				}
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
		}
	}
	for _, sub := range r.Subroutes {
		Bootstrap(sub, group)
	}
}
