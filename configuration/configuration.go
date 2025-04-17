package configuration

import (
	"strconv"
)

type ConfigurationType int8

const (
	RouteNameType            ConfigurationType = iota //ok
	RoutePrefixType                                   //ok
	NetworkCachingPolicyType                          //ok
	SerializationGroupsType                           //ok

	MaxItemPerPageType           //ok
	ItemPerPageType              //ok
	ItemPerPageParameterNameType //ok
	PaginationType               //ok
	PaginationClientControlType  //ok
	PaginationParameterNameType  //ok
	BatchIdsParameterNameType    //ok
	PageParameterNameType        //ok

	SortingClientControlType //ok
	SortingType              //ok
	SortingParameterNameType //ok
	SortableFieldsType       //ok

	//unimplemented
	TypeEnabledType
	DefaultFilteringType
	InMemoryCachingPolicyType
)

type Configuration struct {
	Type   ConfigurationType
	Values []string
}

// default is 0, no caching
// if you set it to 0, it will be disabled
// Be careful with reading policy
func NetworkCachingPolicy(seconds int) Configuration {
	return Configuration{Type: NetworkCachingPolicyType, Values: []string{strconv.Itoa(seconds)}}
}

// default is "api" do not enter / manualy
// for api/v1/ use RoutePrefix("api", "v1")
func RoutePrefix(prefix ...string) Configuration {
	return Configuration{Type: RoutePrefixType, Values: prefix}
}

// by default, it is entity name
func RouteName(name string) Configuration {
	return Configuration{Type: RouteNameType, Values: []string{name}}
}

// serialization groups
func SerializationGroups(groups ...string) Configuration {
	return Configuration{Type: SerializationGroupsType, Values: groups}
}

// default is 1000 per page
func MaxItemPerPage(max int) Configuration {
	return Configuration{Type: MaxItemPerPageType, Values: []string{strconv.Itoa(max)}}
}

// default is 100 per page
func ItemPerPage(defaultValue int) Configuration {
	return Configuration{Type: ItemPerPageType, Values: []string{strconv.Itoa(defaultValue)}}
}

// default is Enabled
// use to enable/disable pagination
// it is recommended to use this option but you might want to disable it
func Pagination(defaultValue bool) Configuration {
	return Configuration{Type: PaginationType, Values: []string{strconv.FormatBool(defaultValue)}}
}

// default is disabled
// allow/disallow client to force pagination using query string
func PaginationClientControl(forced bool) Configuration {
	return Configuration{Type: PaginationClientControlType, Values: []string{strconv.FormatBool(forced)}}
}

// default is "pagination"
// name of the query string parameter used to force pagination
func PaginationParameterName(name string) Configuration {
	return Configuration{Type: PaginationParameterNameType, Values: []string{name}}
}

// default is "page"
func PageParameterName(name string) Configuration {
	return Configuration{Type: PageParameterNameType, Values: []string{name}}
}

// default is "itemsPerPage"
func ItemPerPageParameterName(name string) Configuration {
	return Configuration{Type: ItemPerPageParameterNameType, Values: []string{name}}
}

// default is "ids"
func BatchIdsName(name string) Configuration {
	return Configuration{Type: BatchIdsParameterNameType, Values: []string{name}}
}

// Default is map[string]string{"id": "asc"}
// Converts the map into a slice of strings in the format []string{"id", "asc"}
func Sorting(sortingMap map[string]string) Configuration {
	values := []string{}
	for key, value := range sortingMap {
		values = append(values, key, value)
	}
	return Configuration{Type: SortingType, Values: values}
}

// Default is "sort"
// name of the query string parameter used to sort
func SortingParameterName(name string) Configuration {
	return Configuration{Type: SortingParameterNameType, Values: []string{name}}
}

// Default is true
// allow/disallow client to sort using query string
func SortingClientControl(enabled bool) Configuration {
	return Configuration{Type: SortingClientControlType, Values: []string{strconv.FormatBool(enabled)}}
}

// Default is "id"
// name of the field allowed to be used to sort
func SortableFields(fields ...string) Configuration {
	return Configuration{Type: SortableFieldsType, Values: fields}
}
