package apiman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/errors"
	"github.com/philiphil/apiman/method/MethodType"
	"github.com/philiphil/apiman/router"
	"github.com/philiphil/apiman/serializer/format"
	"strconv"
)

func (r *ApiRouter[T]) GetList(c *gin.Context) {
	pagination, _ := strconv.ParseBool(c.DefaultQuery("pagination", "false"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page--
	itemPerPage, _ := strconv.Atoi(c.DefaultQuery("itemsPerPage", "100"))

	if !pagination {
		itemPerPage = 0
	}

	if pagination {
		objects, err := r.Orm.GetPaginatedList(itemPerPage, page)
		if err != nil {
			panic(err)
		}
		count, err := r.Orm.Count()
		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}
		params := map[string]string{}
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}
		c.Render(
			200,
			router.SerializerRenderer{
				Data:   router.JsonldCollection(objects, c.Request.URL.String(), page+1, params, int(count/int64(itemPerPage))),
				Format: format.JSON,
				Groups: []string{},
			},
		)
		return
	}
	objects, err := r.Orm.GetAll()
	if err != nil {
		c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
	}

	c.Render(200,
		router.SerializerRenderer{
			Data:   objects,
			Format: format.JSON,
			Groups: r.GetMethodConfiguration(method_type.GetList).SerializationGroups,
		})
}
