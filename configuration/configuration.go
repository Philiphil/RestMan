package configuration

import (
	"strconv"
)

// ConfigurationType defines the type of configuration option being set
type ConfigurationType int8

const (
	// RouteNameType sets the route name (default: entity name in snake_case)
	// Example: RouteName("books") generates /api/books
	RouteNameType ConfigurationType = iota

	// RoutePrefixType sets the URL prefix for routes (default: "api")
	// Example: RoutePrefix("api", "v1") generates /api/v1/entity
	RoutePrefixType

	// NetworkCachingPolicyType sets HTTP caching headers (Cache-Control max-age)
	// Value in seconds. Default: 0 (no caching)
	NetworkCachingPolicyType

	// SerializationGroupsType defines which field groups to include in serialization
	// Used with struct tags: `groups:"read,write"`
	SerializationGroupsType

	// MaxItemPerPageType sets the maximum allowed items per page (default: 1000)
	// Prevents clients from requesting too many items at once
	MaxItemPerPageType

	// ItemPerPageType sets the default number of items per page (default: 100)
	ItemPerPageType

	// ItemPerPageParameterNameType sets the query parameter name for items per page (default: "itemsPerPage")
	// Example: ?itemsPerPage=50
	ItemPerPageParameterNameType

	// PaginationType enables or disables pagination (default: enabled)
	// When disabled, all results are returned
	PaginationType

	// PaginationClientControlType allows clients to control pagination via query params (default: disabled)
	// When enabled, clients can use ?page=2&itemsPerPage=50
	PaginationClientControlType

	// PaginationParameterNameType sets the query parameter name for pagination control (default: "pagination")
	PaginationParameterNameType

	// BatchIdsParameterNameType sets the query parameter name for batch operations (default: "ids")
	// Example: GET /api/entity?ids=1,2,3
	BatchIdsParameterNameType

	// PageParameterNameType sets the query parameter name for page number (default: "page")
	// Example: ?page=2
	PageParameterNameType

	// SortingClientControlType allows clients to control sorting via query params (default: enabled)
	// When enabled, clients can use ?order[field]=asc
	SortingClientControlType

	// SortingType sets the default sort order
	// Example: map[string]string{"id": "asc", "createdAt": "desc"}
	SortingType

	// SortingParameterNameType sets the query parameter name for sorting (default: "order")
	// Example: ?order[title]=asc
	SortingParameterNameType

	// SortableFieldsType defines which fields are allowed for sorting (default: "id")
	// Whitelist to prevent sorting on sensitive or non-indexed fields
	SortableFieldsType

	// Unimplemented configuration types - reserved for future use
	BatchLimitType            // Will limit the number of items in batch operations
	TypeEnabledType           // Will enable/disable specific route types
	DefaultFilteringType      // Will add default filters to queries
	InMemoryCachingPolicyType // Will configure in-memory caching
)

// Configuration represents a single configuration option with its type and values.
// Configurations are passed to NewApiRouter to customize router behavior.
//
// Example:
//
//	router := router.NewApiRouter(
//	    orm,
//	    routes,
//	    configuration.ItemPerPage(50),
//	    configuration.MaxItemPerPage(500),
//	)
type Configuration struct {
	Type   ConfigurationType
	Values []string
}

// NetworkCachingPolicy sets the HTTP Cache-Control max-age header in seconds.
// Default is 0 (no caching). Use with caution for frequently changing data.
//
// Example:
//
//	configuration.NetworkCachingPolicy(3600) // Cache for 1 hour
func NetworkCachingPolicy(seconds int) Configuration {
	return Configuration{Type: NetworkCachingPolicyType, Values: []string{strconv.Itoa(seconds)}}
}

// RoutePrefix sets the URL prefix for all routes. Default is "api".
// Do not include leading or trailing slashes.
//
// Example:
//
//	configuration.RoutePrefix("api", "v1") // Generates /api/v1/entity
func RoutePrefix(prefix ...string) Configuration {
	return Configuration{Type: RoutePrefixType, Values: prefix}
}

// RouteName sets the route name. By default, uses the entity name in snake_case.
//
// Example:
//
//	configuration.RouteName("books") // Generates /api/books
func RouteName(name string) Configuration {
	return Configuration{Type: RouteNameType, Values: []string{name}}
}

// SerializationGroups defines which field groups to include in serialization.
// Fields must have matching `groups:"group1,group2"` struct tags.
//
// Example:
//
//	configuration.SerializationGroups("read", "public")
func SerializationGroups(groups ...string) Configuration {
	return Configuration{Type: SerializationGroupsType, Values: groups}
}

// MaxItemPerPage sets the maximum allowed items per page. Default is 1000.
// Prevents clients from requesting excessive data.
//
// Example:
//
//	configuration.MaxItemPerPage(500)
func MaxItemPerPage(max int) Configuration {
	return Configuration{Type: MaxItemPerPageType, Values: []string{strconv.Itoa(max)}}
}

// ItemPerPage sets the default number of items per page. Default is 100.
//
// Example:
//
//	configuration.ItemPerPage(50)
func ItemPerPage(defaultValue int) Configuration {
	return Configuration{Type: ItemPerPageType, Values: []string{strconv.Itoa(defaultValue)}}
}

// Pagination enables or disables pagination. Default is enabled (true).
// When disabled, all results are returned without pagination.
//
// Example:
//
//	configuration.Pagination(false) // Disable pagination
func Pagination(defaultValue bool) Configuration {
	return Configuration{Type: PaginationType, Values: []string{strconv.FormatBool(defaultValue)}}
}

// PaginationClientControl allows clients to control pagination via query parameters.
// Default is disabled (false). When enabled, clients can use ?page=2&itemsPerPage=50.
//
// Example:
//
//	configuration.PaginationClientControl(true)
func PaginationClientControl(forced bool) Configuration {
	return Configuration{Type: PaginationClientControlType, Values: []string{strconv.FormatBool(forced)}}
}

// PaginationParameterName sets the query parameter name for pagination control.
// Default is "pagination".
//
// Example:
//
//	configuration.PaginationParameterName("paginate")
func PaginationParameterName(name string) Configuration {
	return Configuration{Type: PaginationParameterNameType, Values: []string{name}}
}

// PageParameterName sets the query parameter name for page number. Default is "page".
//
// Example:
//
//	configuration.PageParameterName("p") // Use ?p=2 instead of ?page=2
func PageParameterName(name string) Configuration {
	return Configuration{Type: PageParameterNameType, Values: []string{name}}
}

// ItemPerPageParameterName sets the query parameter name for items per page.
// Default is "itemsPerPage".
//
// Example:
//
//	configuration.ItemPerPageParameterName("limit") // Use ?limit=50
func ItemPerPageParameterName(name string) Configuration {
	return Configuration{Type: ItemPerPageParameterNameType, Values: []string{name}}
}

// BatchIdsName sets the query parameter name for batch operations. Default is "ids".
//
// Example:
//
//	configuration.BatchIdsName("id") // Use ?id=1,2,3 instead of ?ids=1,2,3
func BatchIdsName(name string) Configuration {
	return Configuration{Type: BatchIdsParameterNameType, Values: []string{name}}
}

// Sorting sets the default sort order as a map of field names to direction ("asc" or "desc").
// Default is map[string]string{"id": "asc"}.
//
// Example:
//
//	configuration.Sorting(map[string]string{
//	    "createdAt": "desc",
//	    "title": "asc",
//	})
func Sorting(sortingMap map[string]string) Configuration {
	values := []string{}
	for key, value := range sortingMap {
		values = append(values, key, value)
	}
	return Configuration{Type: SortingType, Values: values}
}

// SortingParameterName sets the query parameter name for sorting. Default is "order".
//
// Example:
//
//	configuration.SortingParameterName("sort") // Use ?sort[field]=asc
func SortingParameterName(name string) Configuration {
	return Configuration{Type: SortingParameterNameType, Values: []string{name}}
}

// SortingClientControl allows clients to control sorting via query parameters.
// Default is enabled (true). When enabled, clients can use ?order[field]=asc.
//
// Example:
//
//	configuration.SortingClientControl(false) // Disable client sorting
func SortingClientControl(enabled bool) Configuration {
	return Configuration{Type: SortingClientControlType, Values: []string{strconv.FormatBool(enabled)}}
}

// SortableFields defines which fields are allowed for sorting. Default is ["id"].
// Acts as a whitelist to prevent sorting on sensitive or non-indexed fields.
//
// Example:
//
//	configuration.SortableFields("id", "title", "createdAt")
func SortableFields(fields ...string) Configuration {
	return Configuration{Type: SortableFieldsType, Values: fields}
}
