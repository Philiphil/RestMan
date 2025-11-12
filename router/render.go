package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router/renderer"
	"github.com/philiphil/restman/security"
	"github.com/philiphil/restman/serializer"
)

func (r *ApiRouter[T]) Render(c *gin.Context, statusCode int, obj any, routeType route.RouteType) {
	responseFormat, errParse := ParseAcceptHeader(c.GetHeader("Accept"))
	if errParse != nil {
		abortWithError(c, errParse)
		return
	}
	c.Header("Vary", "Accept, Authorization")

	outputGroups, err := r.GetEffectiveOutputSerializationGroups(c, routeType)
	if err != nil {
		abortWithError(c, err)
		return
	}

	s := serializer.NewSerializer(responseFormat)

	sparseFieldsRequest, err := r.GetSparseFields(c, routeType)
	if err != nil {
		abortWithError(c, err)
		return
	}

	var formatedData string
	var sparseFields map[string][]string
	if sparseFieldsRequest != nil {
		sparseFields = sparseFieldsRequest.ResourceFields
	}
	formatedData, err = s.Serialize(obj, serializer.SerializeOptions{
		Groups:       outputGroups,
		SparseFields: sparseFields,
	})
	if err != nil {
		abortWithError(c, err)
		return
	}
	serializedData := []byte(formatedData)

	etag := GenerateETag(serializedData)
	c.Header("ETag", etag)

	// client catch is tricky
	// if the client provides an If-None-Match header, we must compare it to the generated ETag
	// If they match, it means the client has the latest version and we can respond with 304 Not Modified
	// This helps reduce bandwidth and improve performance by avoiding sending the same data again.

	if match := c.GetHeader("If-None-Match"); match != "" {
		if match == etag {
			c.Status(http.StatusNotModified)
			return // we can stop
		}
	}

	entity := r.Orm.NewEntity()
	visibility := "public"
	if _, hasReadingRights := security.HasReadingRights(entity); hasReadingRights {
		visibility = "private"
	}
	if maxAge, err := r.GetConfiguration(configuration.NetworkCachingPolicyType, routeType); err == nil && len(maxAge.Values) > 0 {
		c.Header("Cache-Control", fmt.Sprintf("%s, max-age=%s", visibility, maxAge.Values[0]))
	}

	c.Data(statusCode, renderer.GetContentType(responseFormat)[0], serializedData)
}
