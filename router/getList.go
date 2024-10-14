package router

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
	method_type "github.com/philiphil/restman/method/MethodType"
)

func (r *ApiRouter[T]) GetList(c *gin.Context) {
	pagination, _ := strconv.ParseBool(c.DefaultQuery("pagination", "false"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page--
	itemPerPage, _ := strconv.Atoi(c.DefaultQuery("itemsPerPage", "100"))

	responseFormat, err := ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	var objects []T

	if pagination {
		objects, err = r.Orm.GetPaginatedList(itemPerPage, page)
		if err != nil {
			panic(err)
		}
		count, err := r.Orm.Count()
		if err != nil {
			panic(err)
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
					Groups: r.GetMethodConfiguration(method_type.GetList).SerializationGroups,
				},
			)
			return
		}
	} else {
		objects, err = r.Orm.GetAll()
		if err != nil {
			c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		}
	}

	c.Render(200,
		SerializerRenderer{
			Data:   objects,
			Format: responseFormat,
			Groups: r.GetMethodConfiguration(method_type.GetList).SerializationGroups,
		})

}
