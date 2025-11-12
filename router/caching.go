package router

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/security"
)

// HandleCaching sets appropriate Cache-Control and Etag headers based on route configuration.
func (r *ApiRouter[T]) HandleCaching(route route.RouteType, c *gin.Context) {
	entity := r.Orm.NewEntity()
	visibility := "public"
	if _, HasReadingRights := security.HasReadingRights(entity); HasReadingRights {
		visibility = "private"
	}

	//todo add error handling
	maxAge, _ := r.GetConfiguration(configuration.NetworkCachingPolicyType, route)
	bodyBytes, _ := io.ReadAll(c.Request.Body)

	c.Header("Etag", GenerateETag(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	c.Header("Cache-Control", visibility+", max-age="+(maxAge.Values[0]))

}

func GenerateETag(responseBody []byte) string {
	hash := sha256.Sum256(responseBody)
	return fmt.Sprintf(`"%x"`, hash[:8])
}
