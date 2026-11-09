package router

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/route"
)

func (r *ApiRouter[T]) GetIncludedRelationships(c *gin.Context, routeType route.RouteType) (results []string) {
	include, err := r.GetConfiguration(configuration.IncludeEnabledType, routeType)
	if err != nil {
		return results
	}
	includeBool, err := strconv.ParseBool(include.Values[0])
	if err != nil || !includeBool {
		return results
	}
	includeParameter, err := r.GetConfiguration(configuration.IncludeParameterNameType, routeType)
	if err != nil {
		return results
	}
	includeQuery := c.Query(includeParameter.Values[0])
	if includeQuery == "" {
		return results
	}
	results = strings.Split(includeQuery, ",")

	return results
}

// Get handles HTTP GET requests to retrieve a single entity by ID.
func (r *ApiRouter[T]) Get(c *gin.Context) {
	include := r.GetIncludedRelationships(c, route.Get)

	repo, err := r.Orm.Repo.LoadAssociations(include)
	if err != nil {
		abortWithError(c, errors.Internal)
		return
	}

	object, err := orm.NewORM(repo).GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		abortWithApiError(c, errors.NotFound)
		return
	}

	err = r.ReadingCheck(c, object)
	if err != nil {
		abortWithError(c, err)
		return
	}
	fmt.Printf("Retrieved object: %+v\n", object)
	r.Render(c, 200, object, route.Get)

}
