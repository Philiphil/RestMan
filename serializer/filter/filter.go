package filter

import (
	"reflect"
)

func FilterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if isStruct(elemType) {
		return filterByGroupsStruct(obj, groups...)
	}
	if isList(elemType) {
		return filterByGroupsSlice(obj, groups...)
	}
	if isMap(elemType) {
		return filterByGroupsMap(obj, groups...)
	}
	return obj
}
