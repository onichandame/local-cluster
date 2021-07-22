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
	"github.com/sirupsen/logrus"
)

var Entrances = route.Route{
	Endpoint: "/entrances/:name/*proxyPath",
	Handler: func(c *gin.Context) (interface{}, error) {
		entrance := model.Entrance{}
		if err := db.Db.Preload("Backend").Where("name = ?", c.Param("name")).First(&entrance).Error; err != nil {
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
		proxy.ModifyResponse = func(r *http.Response) error {
			r.Request.Host = c.Request.Host
			r.Request.URL.Host = c.Request.URL.Host
			r.Request.URL.Path = c.Request.URL.Path
			r.Request.URL.Scheme = c.Request.URL.Scheme
			return nil
		}
		proxy.ServeHTTP(c.Writer, c.Request)
		logrus.Warn("request finished")
		return nil, nil
	},
}
