package ApiRouter

import (
	"ApiMan/Gin"
	"ApiMan/Serializer/Format"
	"github.com/gin-gonic/gin"
	"strconv"
)

func (r *ApiRouter[T]) GetList(c *gin.Context) {
	pagination, _ := strconv.ParseBool(c.DefaultQuery("pagination", "false"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page--
	itemPerPage, _ := strconv.Atoi(c.DefaultQuery("itemPerPage", "100"))

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
			Gin.SerializerRenderer{
				Data:   Gin.JsonldCollection(objects, c.Request.URL.String(), page+1, params, int(count/int64(itemPerPage))),
				Format: Format.JSON,
				Groups: []string{},
			},
		)
		return
	}
	objects, err := r.Orm.GetAll()
	if err != nil {
		panic(err)
	}

	c.Render(200,
		Gin.SerializerRenderer{
			Data:   objects,
			Format: Format.JSON,
			Groups: []string{},
		})
}
