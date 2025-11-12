package filter

import (
	"reflect"
)

func filterByGroupsSlice[T any](obj T, groups []string, sparseFields map[string][]string, currentPath string) T {
	value := DereferenceValueIfPointer(reflect.ValueOf(obj))

	if value.Len() == 0 {
		return obj
	}

	firstElem := value.Index(0).Interface()

	filteredFirstElem := filterByGroupsWithPath(firstElem, groups, sparseFields, currentPath)

	newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(filteredFirstElem)), 0, value.Len())
	newSlice = reflect.Append(newSlice, reflect.ValueOf(filteredFirstElem))

	for i := 1; i < value.Len(); i++ {
		elem := value.Index(i)
		filteredElem := filterByGroupsWithPath(elem.Interface(), groups, sparseFields, currentPath)
		newSlice = reflect.Append(newSlice, reflect.ValueOf(filteredElem))
	}
	return newSlice.Interface().(T)
}
