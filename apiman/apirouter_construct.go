package apiman

import (
	"github.com/philiphil/apiman/method"
	"github.com/philiphil/apiman/orm"
	"github.com/philiphil/apiman/orm/entity"
	"reflect"
	"strings"
)

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
