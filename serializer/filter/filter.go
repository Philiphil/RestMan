package filter

import (
	"reflect"
)

// FilterByGroups filters an object's fields based on the provided serialization groups.
// If sparseFields is provided, it will further filter fields to only include those specified.
// sparseFields is a map where keys are resource paths ("" for root, "posts" for nested) and values are field names.
func FilterByGroups[T any](obj T, groups []string, sparseFields map[string][]string) T {
	return filterByGroupsWithPath(obj, groups, sparseFields, "")
}

func filterByGroupsWithPath[T any](obj T, groups []string, sparseFields map[string][]string, currentPath string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if IsStruct(elemType) {
		return filterByGroupsStruct(obj, groups, sparseFields, currentPath)
	}
	if IsList(elemType) {
		return filterByGroupsSlice(obj, groups, sparseFields, currentPath)
	}
	if IsMap(elemType) {
		return filterByGroupsMap(obj, groups, sparseFields, currentPath)
	}
	return obj
}
