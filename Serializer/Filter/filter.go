package Filter

import (
	"reflect"
)

func FilterByGroups[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if isStruct(elemType) {
		return filterByGroups_struct(obj, groups...)
	}
	if isList(elemType) {
		return filterByGroups_slice(obj, groups...)
	}
	if isMap(elemType) {
		return filterByGroups_map(obj, groups...)
	}
	return obj
}
