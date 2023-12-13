package ApiMan

import (
	"github.com/gin-gonic/gin"
	"github.com/philiphil/apiman/method"
	method_type "github.com/philiphil/apiman/method/MethodType"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
	"reflect"
	"strings"
	"unicode"
)

type ApiRouter[T entity.IEntity] struct {
	Orm     orm.ORM[T]
	Methods []method.ApiMethodConfiguration
	//AuthProvider
	Route string
}

func (r *ApiRouter[T]) AllowRoutes(router *gin.Engine) {
	for _, method_ := range r.Methods {
		switch method_.Method {
		case method_type.Get:
			router.GET(r.Route+"/:id", r.Get)
			router.HEAD(r.Route+"/:id", r.Head)
		case method_type.GetList:
			router.GET(r.Route, r.GetList)
		case method_type.Post:
			router.POST(r.Route, r.Post)
		case method_type.Put:
			router.PUT(r.Route+"/:id", r.Put)
		case method_type.Patch:
			router.PATCH(r.Route+"/:id", r.Patch)
		case method_type.Delete:
			router.DELETE(r.Route+"/:id", r.Delete)
		case method_type.Undefined:
		case method_type.Connect:
		case method_type.Trace:
		case method_type.Options:
		}
	}
	return
}

func NewApiRouter[T entity.IEntity](orm orm.ORM[T], methods []method.ApiMethodConfiguration, route ...string) *ApiRouter[T] {
	router := &ApiRouter[T]{
		Orm:     orm,
		Methods: methods,
	}
	if len(route) > 0 {
		if !strings.HasPrefix(route[0], "/") {
			router.Route = "/" + route[0]
		} else {
			router.Route = route[0]
		}
		if !strings.HasSuffix(router.Route, "/") {
			router.Route = strings.TrimSuffix(router.Route, "/")
		}
	} else {
		router.Route = "/api/" + convertToSnakeCase(reflect.TypeOf(orm.NewEntity()).Name())
	}

	return router

}

func convertToSnakeCase(input string) string {
	runes := []rune(input)
	if len(runes) == 0 {
		return ""
	}

	runes[0] = unicode.ToLower(runes[0])

	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) {
			runes[i] = unicode.ToLower(runes[i])
			runes = append(runes[:i], append([]rune{'_'}, runes[i:]...)...)
			i++
		}
	}

	return string(runes)
}
func (r ApiRouter[T]) GetMethodConfiguration(apiMethod method_type.ApiMethod) method.ApiMethodConfiguration {
	for _, method_ := range r.Methods {
		if method_.Method == apiMethod {
			return method_
		}
	}
	return method.New()
}
