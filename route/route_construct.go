package route

import (
	"maps"

	"github.com/philiphil/restman/configuration"
)

func NewRoute(routeType RouteType, configurations ...configuration.Configuration) Route {
	c := Route{}
	c.RouteType = routeType
	for _, configuration := range configurations {
		c.Configuration[configuration.Type] = configuration
	}
	return c
}

// DefaultApiRoutes returns a map of default routes
// Get, GetList, Post, Put, Patch, Delete, Head, Options
// with empty configurations
// This is useful for creating a default ApiRouter with the standards operations
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

func AllApiRoutes() map[RouteType]Route {
	mergedRoutes := DefaultApiRoutes()
	maps.Copy(mergedRoutes, BatchOperations())
	return mergedRoutes
}

// BatchOperations
func BatchOperations() map[RouteType]Route {
	return map[RouteType]Route{
		BatchDelete: NewRoute(BatchDelete),
		BatchPut:    NewRoute(BatchPut),
		BatchPatch:  NewRoute(BatchPatch),
		BatchPost:   NewRoute(BatchPost),
		BatchGet:    NewRoute(BatchGet),
	}
}
