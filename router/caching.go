package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/security"
)

func (r *ApiRouter[T]) HandleCaching(route route.RouteType, c *gin.Context) {
	entity := r.Orm.NewEntity()
	visibility := "public"
	if _, HasReadingRights := security.HasReadingRights(entity); HasReadingRights {
		visibility = "private"
	}

	maxAge, err := r.GetConfiguration(configuration.NetworkCachingPolicyType, route)

	if maxAge.Values[0] != "0" && err == nil {
		c.Header("Cache-Control", visibility+", max-age="+(maxAge.Values[0]))
		c.Header("Etag", c.Request.URL.String())
	}
}
