package ApiMan

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/method"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
)

type ApiRouter[T entity.IEntity] struct {
	Orm     orm.ORM[T]
	Methods []method.ApiMethod
	Prefix  string
	Name    string
}

func (r *ApiRouter[T]) AllowRoutes(router *gin.Engine) {
	for _, method_ := range r.Methods {
		switch method_ {
		case method.Get:
			router.GET("/"+r.Prefix+"/"+r.Name+"/:id", r.Get)
			router.HEAD("/"+r.Prefix+"/"+r.Name+"/:id", r.Head)
		case method.GetList:
			router.GET("/"+r.Prefix+"/"+r.Name, r.GetList)
		case method.Post:
			router.POST("/"+r.Prefix+"/"+r.Name, r.Post)
		case method.Put:
			router.PUT("/"+r.Prefix+"/"+r.Name+"/:id", r.Put)
		case method.Patch:
			router.PATCH("/"+r.Prefix+"/"+r.Name+"/:id", r.Patch)
		case method.Delete:
			router.DELETE("/"+r.Prefix+"/"+r.Name+"/:id", r.Delete)
		case method.Undefined:
		case method.Connect:
		case method.Trace:
		case method.Options:
		}
	}
	return
}

func NewApiRouter[T entity.IEntity](orm orm.ORM[T], methods []method.ApiMethod, prefix, name string) *ApiRouter[T] {
	return &ApiRouter[T]{
		Orm:     orm,
		Methods: methods,
		Prefix:  prefix,
		Name:    name,
	}
}
