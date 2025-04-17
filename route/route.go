package route

import "github.com/philiphil/restman/configuration"

// Route is a struct that represents a route
// it use a type RouteType to define the type of route (Get, Post, Put, etc)
// and a map of configuration.Configuration to store the configurations of the route
type Route struct {
	RouteType     RouteType
	Configuration map[configuration.ConfigurationType]configuration.Configuration
}
