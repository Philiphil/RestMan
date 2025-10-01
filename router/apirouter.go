package router

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/entity"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/security"
)

// SubresourceRegistrar is an interface that any ApiRouter must implement
// to allow registration of its routes with a parent router
type SubresourceRegistrar interface {
	// RegisterSubroutes registers all routes for this subresource under the given parent route
	RegisterSubroutes(router *gin.Engine, parentRoute string)
	// GetSubresourceName returns the name of this subresource (used in URL path)
	GetSubresourceName() string
}

// An ApiRouter is the main object to create a REST API
// It is composed of an ORM, a list of Allow methods, a list of firewalls and a route
// To create an ApiRouter, you should use the NewApiRouter function
type ApiRouter[T entity.Entity] struct {
	Orm       orm.ORM[T]
	Routes    map[route.RouteType]route.Route
	Firewalls []security.Firewall

	Configuration map[configuration.ConfigurationType]configuration.Configuration

	Subresources []SubresourceRegistrar
}

// AllowRoutes is a function that adds the route to the gin router
func (r *ApiRouter[T]) AllowRoutes(router *gin.Engine) {

	//Batch Get and Bast Post shares the same route as GetList and Post
	//we dont want to register the route twice
	getList, post := false, false

	for _, route_ := range r.Routes {
		routeName := r.Route(route_.RouteType)
		switch route_.RouteType {
		case route.Get:
			router.GET(routeName+"/:id", r.Get)
		case route.BatchGet, route.GetList:
			if !getList {
				router.GET(routeName, r.GetListOrBatchGet)
				getList = true
			}
		case route.BatchPost, route.Post:
			if !post {
				router.POST(routeName, r.Post)
				post = true
			}
		case route.Put:
			router.PUT(routeName+"/:id", r.Put)
		case route.Patch:
			router.PATCH(routeName+"/:id", r.Patch)
		case route.Delete:
			router.DELETE(routeName+"/:id", r.Delete)
		case route.Head:
			router.HEAD(routeName+"/:id", r.Head)
		case route.Options:
			router.OPTIONS(routeName+"/:id", r.Options)
			router.OPTIONS(routeName, r.Options)
		case route.BatchDelete:
			router.DELETE(routeName, r.batchDelete)
		case route.BatchPatch:
			router.PATCH(routeName, r.BatchPatch)
		case route.BatchPut:
			router.PUT(routeName, r.BatchPut)
		case route.Connect:
		case route.Trace:
		case route.Undefined:
		}
	}

	// Register all subresources
	baseroute := r.Route()
	for _, subresource := range r.Subresources {
		subresource.RegisterSubroutes(router, baseroute)
	}
}

func ConvertToSnakeCase(input string) string {
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

// This function return either the router wide configuration or the route specific configuration
// If the routeType is not provided, it will return the router wide configuration
// If the routeType is provided, it will return the route specific configuration
// error is returned if the configuration is not found
// by default error should always be nil if you use NewApiRouter
func (r *ApiRouter[T]) GetConfiguration(configuration configuration.ConfigurationType, routeType ...route.RouteType) (configuration.Configuration, error) {
	routerValue, found := r.Configuration[configuration]
	if len(routeType) == 1 {
		for _, route_ := range r.Routes {
			if route_.RouteType == routeType[0] {
				routeValue, exists := route_.Configuration[configuration]
				if exists {
					return routeValue, nil
				}
			}
		}
	}
	if !found {
		return routerValue, errors.ApiError{Code: errors.ErrInternal.Code, Message: errors.ErrInternal.Message}
	}

	return routerValue, nil
}

// NewApiRouter is a function that creates a new ApiRouter
// it should be the default way of creating an ApiRouter because it sets the default configuration
func NewApiRouter[T entity.Entity](orm orm.ORM[T], routes map[route.RouteType]route.Route, conf ...configuration.Configuration) *ApiRouter[T] {
	router := &ApiRouter[T]{
		Orm:    orm,
		Routes: routes,
	}

	router.Configuration = configuration.DefaultConfiguration()
	routeNameSet := false
	for _, confV := range conf {
		if confV.Type == configuration.RouteNameType {
			routeNameSet = true
		}
		router.Configuration[confV.Type] = confV
	}
	//The default RouteName is the name of the entity in snake case
	//and it cannot be decided in advance so it is not set by configuration.DefaultConfiguration()
	//it should be the only one configuration without a default value
	if !routeNameSet {
		router.Configuration[configuration.RouteNameType] = configuration.RouteName(ConvertToSnakeCase(reflect.TypeOf(orm.NewEntity()).Name()))
	}
	return router
}

func TrimSlash(s string) string {
	return strings.TrimSuffix(strings.TrimPrefix(s, "/"), "/")
}

// Route is a function that returns the route name for a given route type
func (r *ApiRouter[T]) Route(routeType ...route.RouteType) (name string) {
	name = "/"
	prefixs, _ := r.GetConfiguration(configuration.RoutePrefixType, routeType...)
	for _, v := range prefixs.Values {
		name += TrimSlash(v) + "/"
	}

	routeName, _ := r.GetConfiguration(configuration.RouteNameType, routeType...)

	name += TrimSlash(routeName.Values[0])
	return name
}

func (r *ApiRouter[T]) AddFirewall(firewall ...security.Firewall) {
	r.Firewalls = append(r.Firewalls, firewall...)
}

// AddSubresource adds a subresource to this ApiRouter
// The subresource routes will be registered under the parent route with /:id/ prefix
func (r *ApiRouter[T]) AddSubresource(subresource SubresourceRegistrar) {
	r.Subresources = append(r.Subresources, subresource)
}

// GetSubresourceName returns the name of this resource (used when registered as a subresource)
func (r *ApiRouter[T]) GetSubresourceName() string {
	routeName, _ := r.GetConfiguration(configuration.RouteNameType)
	return routeName.Values[0]
}

// RegisterSubroutes registers all routes for this ApiRouter as a subresource under the given parent route
func (r *ApiRouter[T]) RegisterSubroutes(router *gin.Engine, parentRoute string) {
	//Batch Get and Bast Post shares the same route as GetList and Post
	//we dont want to register the route twice
	getList, post := false, false

	subresourceName := r.GetSubresourceName()

	// Always use "id" for all parameters to avoid gin routing conflicts
	paramName := "id"
	itemParamName := "id"

	baseRoute := parentRoute + "/:" + paramName + "/" + subresourceName

	for _, route_ := range r.Routes {
		switch route_.RouteType {
		case route.Get:
			router.GET(baseRoute+"/:"+itemParamName, r.Get)
		case route.BatchGet, route.GetList:
			if !getList {
				router.GET(baseRoute, r.GetListOrBatchGet)
				getList = true
			}
		case route.BatchPost, route.Post:
			if !post {
				router.POST(baseRoute, r.Post)
				post = true
			}
		case route.Put:
			router.PUT(baseRoute+"/:"+itemParamName, r.Put)
		case route.Patch:
			router.PATCH(baseRoute+"/:"+itemParamName, r.Patch)
		case route.Delete:
			router.DELETE(baseRoute+"/:"+itemParamName, r.Delete)
		case route.Head:
			router.HEAD(baseRoute+"/:"+itemParamName, r.Head)
		case route.Options:
			router.OPTIONS(baseRoute+"/:"+itemParamName, r.Options)
			router.OPTIONS(baseRoute, r.Options)
		case route.BatchDelete:
			router.DELETE(baseRoute, r.batchDelete)
		case route.BatchPatch:
			router.PATCH(baseRoute, r.BatchPatch)
		case route.BatchPut:
			router.PUT(baseRoute, r.BatchPut)
		case route.Connect:
		case route.Trace:
		case route.Undefined:
		}
	}

	// Recursively register nested subresources
	for _, sub := range r.Subresources {
		sub.RegisterSubroutes(router, baseRoute)
	}
}
