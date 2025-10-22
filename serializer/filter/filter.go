package filter

import (
	"reflect"
)

// FilterByGroups filters an object's fields based on the provided serialization groups.
func FilterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if IsStruct(elemType) {
		return filterByGroupsStruct(obj, groups...)
	}
	if IsList(elemType) {
		return filterByGroupsSlice(obj, groups...)
	}
	if IsMap(elemType) {
		return filterByGroupsMap(obj, groups...)
	}
	return obj
}
