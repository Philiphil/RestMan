package router

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router/renderer"
)

// GetList handles HTTP GET requests to retrieve a collection of entities with optional pagination and sorting.
func (r *ApiRouter[T]) GetList(c *gin.Context) {
	paginate, err := r.IsPaginationEnabled(c)
	if err != nil {
		abortWithError(c, err)
		return
	}
	itemPerPage, err := r.GetItemPerPage(c)
	if err != nil {
		abortWithError(c, err)
		return
	}
	sortOrder, err := r.GetSortOrder(c)
	if err != nil {
		abortWithError(c, err)
		return
	}
	page, err := r.GetPage(c)
	if err != nil {
		abortWithError(c, err)
		return
	}

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		abortWithError(c, err)
		return
	}

	groups, err := r.GetEffectiveOutputSerializationGroups(c, route.GetList)
	if err != nil {
		abortWithApiError(c, errors.Internal)
		return
	}

	var objects []T
	if paginate {
		objects, err = r.Orm.GetPaginatedList(c.Request.Context(), itemPerPage, page, sortOrder)
		if err != nil {
			abortWithApiError(c, errors.DatabaseIssue)
			return
		}
		count, err := r.Orm.Count(c.Request.Context())
		if err != nil {
			abortWithApiError(c, errors.DatabaseIssue)
			return
		}

		params := map[string]string{}
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		if responseFormat == format.JSONLD {
			// Pass the total count (not just current page size) for proper Hydra compliance
			c.Render(
				200,
				renderer.SerializerRenderer{
					Data:   JsonldCollection(objects, int(count), c.Request.URL.String(), page+1, params, int((count+int64(itemPerPage)-1)/int64(itemPerPage))),
					Format: responseFormat,
					Groups: groups,
				},
			)
			return
		}
	} else {
		objects, err = r.Orm.GetAll(c.Request.Context(), sortOrder)
		if err != nil {
			abortWithApiError(c, errors.DatabaseIssue)
			return
		}
	}

	r.Render(c, 200, objects, route.GetList)

}
