package ApiRouter

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/ApiRouter/Method"
	"github.com/philiphil/apiman/Gorm"
)

type ApiRouter[T Gorm.IEntity] struct {
	Orm     Gorm.ORM[T]
	Methods []Method.ApiMethod
	Prefix  string
	Name    string
}

func (r *ApiRouter[T]) AllowRoutes(router *gin.Engine) {
	for _, method := range r.Methods {
		switch method {
		case Method.Get:
			router.GET("/"+r.Prefix+"/"+r.Name+"/:id", r.Get)
			router.HEAD("/"+r.Prefix+"/"+r.Name+"/:id", r.Head)
		case Method.GetList:
			router.GET("/"+r.Prefix+"/"+r.Name, r.GetList)
		case Method.Post:
			router.POST("/"+r.Prefix+"/"+r.Name, r.Post)
		case Method.Put:
			router.PUT("/"+r.Prefix+"/"+r.Name+"/:id", r.Put)
		case Method.Patch:
			router.PATCH("/"+r.Prefix+"/"+r.Name+"/:id", r.Patch)
		case Method.Delete:
			router.DELETE("/"+r.Prefix+"/"+r.Name+"/:id", r.Delete)
		case Method.Undefined:
		case Method.Connect:
		case Method.Trace:
		case Method.Options:
		}
	}
	return
}

func NewApiRouter[T Gorm.IEntity](orm Gorm.ORM[T], methods []Method.ApiMethod, prefix, name string) *ApiRouter[T] {
	return &ApiRouter[T]{
		Orm:     orm,
		Methods: methods,
		Prefix:  prefix,
		Name:    name,
	}
}
