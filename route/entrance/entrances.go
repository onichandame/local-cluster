package entrance

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/onichandame/local-cluster/db"
	"github.com/onichandame/local-cluster/db/model"
	"github.com/onichandame/local-cluster/pkg/route"
)

var Entrances = route.Route{
	Endpoint: "/entrances/:id/*proxyPath",
	Handler: func(c *gin.Context) (interface{}, error) {
		entrance := model.Entrance{}
		if err := db.Db.Preload("Backend").First(&entrance, c.Param("id")).Error; err != nil {
			return nil, err
		}
		uri, err := url.Parse(fmt.Sprintf("http://localhost:%d", entrance.Backend.Port))
		if err != nil {
			return nil, err
		}
		proxy := httputil.NewSingleHostReverseProxy(uri)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = uri.Host
			req.URL.Scheme = uri.Scheme
			req.URL.Host = uri.Host
			req.URL.Path = c.Param("proxyPath")
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		return nil, nil
	},
}
