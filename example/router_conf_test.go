package example_test

import (
	"testing"

	"github.com/philiphil/restman/configuration"
	"github.com/philiphil/restman/orm"
	"github.com/philiphil/restman/orm/gormrepository"
	"github.com/philiphil/restman/route"
	"github.com/philiphil/restman/router"
	"golang.org/x/exp/maps"
)

// As an entity we will rely on the Test struct define in the basic_router_test.go file
// In this file we'll demonstrate how we can customize the router configuration

func TestRouterConfiguration(t *testing.T) {
	//when we do this
	api := router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
	)
	//in reality we are doing this
	api = router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
		maps.Values(configuration.DefaultConfiguration())...,
	)
	//it is why its really important to only create ApiRouter using the NewApiRouter function
	//otherwise you will have to fill yourself all the configuration in order to run

	//this is what configuration.DefaultConfiguration()) is really doing
	api = router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
		configuration.RoutePrefix("api"),
		configuration.NetworkCachingPolicy(0),
		configuration.InputSerializationGroups(),
		configuration.Pagination(true),
		configuration.PaginationClientControl(false),
		configuration.ItemPerPage(100),
		configuration.MaxItemPerPage(1000),
	)

	//when you want to change something like force the pagination
	api = router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		route.DefaultApiRoutes(),
		configuration.PaginationClientControl(true),
	)
	//your  configuration.PaginationClientControl(true) will actually override the default forced pagination value defined in configuration.DefaultConfiguration()

	//You can also set route specific configuration, for example for route
	routes := route.DefaultApiRoutes()
	routes[route.Get].Configuration[configuration.NetworkCachingPolicyType] = configuration.NetworkCachingPolicy(1)
	api = router.NewApiRouter(
		*orm.NewORM(gormrepository.NewRepository[Test](getDB())),
		routes,
	)
	//this way the configuration.NetworkCachingPolicy will be overriden only for this route

	api.AllowRoutes(SetupRouter())
}
