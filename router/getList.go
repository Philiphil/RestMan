package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/route"
)

// GetList handles HTTP GET requests to retrieve a collection of entities with optional pagination and sorting.
func (r *ApiRouter[T]) GetList(c *gin.Context) {
	paginate, err := r.IsPaginationEnabled(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	itemPerPage, err := r.GetItemPerPage(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	sortOrder, err := r.GetSortOrder(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}
	page, err := r.GetPage(c)
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	groups, err := r.GetEffectiveOutputSerializationGroups(c, route.GetList)
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrInternal.Code, errors.ErrInternal.Message)
		return
	}

	var objects []T
	if paginate {
		objects, err = r.Orm.GetPaginatedList(itemPerPage, page, sortOrder)
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
			return
		}
		count, err := r.Orm.Count()
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
			return
		}
		params := map[string]string{}
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		if responseFormat == format.JSONLD {
			c.Render(
				200,
				SerializerRenderer{
					Data:   JsonldCollection(objects, c.Request.URL.String(), page+1, params, int((count+int64(itemPerPage)-1)/int64(itemPerPage))),
					Format: responseFormat,
					Groups: groups,
				},
			)
			return
		}
	} else {
		objects, err = r.Orm.GetAll(sortOrder)
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		}
	}

	c.Render(200,
		SerializerRenderer{
			Data:   objects,
			Format: responseFormat,
			Groups: groups,
		})

}
