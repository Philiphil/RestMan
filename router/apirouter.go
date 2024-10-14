package router

import (
	"reflect"
	"strings"

	"github.com/philiphil/restman/method"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/security"
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
		router.Route = strings.TrimSuffix(router.Route, "/")

	} else {
		router.Route = "/api/" + ConvertToSnakeCase(reflect.TypeOf(orm.NewEntity()).Name())
	}
	return router
}

func (r *ApiRouter[T]) AddFirewall(firewall ...security.Firewall) {
	r.Firewalls = append(r.Firewalls, firewall...)
}
