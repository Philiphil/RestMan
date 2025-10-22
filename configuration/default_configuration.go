package configuration

// DefaultConfiguration returns the default configuration map used by an ApiRouter.
func DefaultConfiguration() map[ConfigurationType]Configuration {
	return map[ConfigurationType]Configuration{
		RoutePrefixType:              RoutePrefix("api"),
		NetworkCachingPolicyType:     NetworkCachingPolicy(0),
		SerializationGroupsType:      SerializationGroups(),
		PaginationType:               Pagination(true),
		PageParameterNameType:        PageParameterName("page"),
		PaginationClientControlType:  PaginationClientControl(false),
		PaginationParameterNameType:  PaginationParameterName("pagination"),
		ItemPerPageType:              ItemPerPage(100),
		MaxItemPerPageType:           MaxItemPerPage(1000),
		BatchIdsParameterNameType:    BatchIdsName("ids"),
		ItemPerPageParameterNameType: ItemPerPageParameterName("itemsPerPage"),

		SortingClientControlType: SortingClientControl(true),
		SortingType:              Sorting(map[string]string{"id": "asc"}),
		SortingParameterNameType: SortingParameterName("sort"),
		SortableFieldsType:       SortableFields("id"),
	}
}
