package router

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
)

// what is the difference between BatchGet and GetList?
// BatchGet is a route that returns a list of objects using given ids
// GetList is a route that returns a list of objects using pagination
func (r *ApiRouter[T]) IsBatchGetOrGetList(c *gin.Context) route.RouteType {
	//first if BatchGet or GetList is not allowed, it is not a BatchGet or GetList
	if _, ok := r.Routes[route.BatchGet]; !ok {
		return route.GetList
	}
	if _, ok := r.Routes[route.GetList]; !ok {
		return route.BatchGet
	}

	ids, _ := r.GetConfiguration(configuration.BatchIdsParameterNameType, route.BatchGet)
	idsParameter := ids.Values[0]
	exists := false
	if _, exists = c.GetQuery(idsParameter); !exists {
		_, exists = c.GetQuery(idsParameter + "[]")
	}
	if exists {
		return route.BatchGet
	}
	return route.GetList
}

func (r *ApiRouter[T]) GetListOrBatchGet(c *gin.Context) {
	rr := r.IsBatchGetOrGetList(c)
	if rr == route.BatchGet {
		r.BatchGet(c)
	} else {
		r.GetList(c)
	}
}

func (r *ApiRouter[T]) GetIds(c *gin.Context) []string {
	ids, _ := r.GetConfiguration(configuration.BatchIdsParameterNameType, route.BatchGet)
	idsParameter := ids.Values[0]
	exists := false
	if _, exists = c.GetQuery(idsParameter + "[]"); exists {
		return c.QueryArray(idsParameter + "[]")
	}

	idsValues := c.QueryArray(ids.Values[0])
	if len(idsValues) == 1 && len(strings.Split(idsValues[0], ",")) > 1 {
		return strings.Split(idsValues[0], ",")
	}
	return idsValues
}

func (r *ApiRouter[T]) BatchGet(c *gin.Context) {
	idsValues := r.GetIds(c)

	formatedId := make([]entity.ID, len(idsValues))
	for i, v := range idsValues {
		formatedId[i] = entity.CastId(v)
	}
	objects, err := r.Orm.FindByIDs(formatedId)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrNotFound.Code, errors.ErrNotFound.Message)
		return
	}
	for _, object := range objects {
		err = r.ReadingCheck(c, object)
		if err != nil {
			c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
			return
		}
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetConfiguration(configuration.SerializationGroupsType, route.Get)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	c.Render(200, SerializerRenderer{
		Data:   objects,
		Format: responseFormat,
		Groups: groups.Values,
	})
}
