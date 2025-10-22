package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/route"
)

// Options handles HTTP OPTIONS requests to return allowed HTTP methods for the resource.
func (r *ApiRouter[T]) Options(c *gin.Context) {
	allowed := ""
	for _, method := range r.Routes {
		switch method.RouteType {
		case route.GetList:
		case route.BatchDelete:
		case route.BatchPatch:
		case route.BatchPut:
		case route.BatchGet:
		case route.BatchPost:
		default:
			if len(allowed) > 0 {
				allowed += ","
			}
			allowed += method.RouteType.String()
		}
	}
	c.Header("Allow", allowed)
	c.Header("Content-Length", "0")
	c.Status(200)
}
