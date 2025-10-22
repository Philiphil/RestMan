package route

import (
	"maps"

	"github.com/philiphil/restman/configuration"
)

// NewRoute creates a new Route with the specified route type and optional configurations.
func NewRoute(routeType RouteType, configurations ...configuration.Configuration) Route {
	c := Route{}
	c.RouteType = routeType
	for _, configuration := range configurations {
		c.Configuration[configuration.Type] = configuration
	}
	return c
}

// DefaultApiRoutes returns a map of default CRUD routes (Get, GetList, Post, Put, Patch, Delete, Head, Options) with empty configurations.
func DefaultApiRoutes() map[RouteType]Route {
	return map[RouteType]Route{
		Get:     NewRoute(Get),
		GetList: NewRoute(GetList),
		Post:    NewRoute(Post),
		Put:     NewRoute(Put),
		Patch:   NewRoute(Patch),
		Delete:  NewRoute(Delete),
		Head:    NewRoute(Head),
		Options: NewRoute(Options),
	}
}

// AllApiRoutes returns a map containing all default routes and batch operation routes.
func AllApiRoutes() map[RouteType]Route {
	mergedRoutes := DefaultApiRoutes()
	maps.Copy(mergedRoutes, BatchOperations())
	return mergedRoutes
}

// BatchOperations returns a map of batch operation routes (BatchDelete, BatchPut, BatchPatch, BatchPost, BatchGet).
func BatchOperations() map[RouteType]Route {
	return map[RouteType]Route{
		BatchDelete: NewRoute(BatchDelete),
		BatchPut:    NewRoute(BatchPut),
		BatchPatch:  NewRoute(BatchPatch),
		BatchPost:   NewRoute(BatchPost),
		BatchGet:    NewRoute(BatchGet),
	}
}
