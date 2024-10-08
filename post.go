package restman

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/errors"
	method_type "github.com/philiphil/restman/method/MethodType"
	"github.com/philiphil/restman/router"
)

func (r *ApiRouter[T]) Post(c *gin.Context) {
	entity := r.Orm.NewEntity()
	if err := r.WritingCheck(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if err := router.UnserializeBodyAndMerge(c, &entity); err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	if err := r.Orm.Create(&entity); err != nil {
		c.AbortWithStatusJSON(errors.ErrDatabaseIssue.Code, errors.ErrDatabaseIssue.Message)
		return
	}
	responseFormat, err := router.ParseAcceptHeader(c.GetHeader("Accept"))
	if err != nil {
		c.AbortWithStatusJSON(err.(errors.ApiError).Code, err.(errors.ApiError).Message)
		return
	}

	c.Render(201, router.SerializerRenderer{
		Data:   &entity,
		Format: responseFormat,
		Groups: r.GetMethodConfiguration(method_type.Get).SerializationGroups, //what is sent shall be compliant to get
	})
}
